package probe

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Alaxay8/routeflux/internal/domain"
)

// TCPChecker checks node health using a TCP connect probe.
type TCPChecker struct {
	Timeout time.Duration
	Now     func() time.Time
}

// Check probes a node endpoint and measures TCP connect latency.
func (c TCPChecker) Check(ctx context.Context, node domain.Node) Result {
	timeout := c.Timeout
	if timeout <= 0 {
		timeout = 5 * time.Second
	}

	now := c.Now
	if now == nil {
		now = time.Now
	}

	if node.Protocol == domain.ProtocolHysteria || node.Protocol == domain.ProtocolHysteria2 {
		latency, err := pingICMP(ctx, node.Address, timeout)
		if err != nil {
			if isCommandNotFound(err) {
				return Result{
					NodeID:  node.ID,
					Healthy: true,
					Checked: now(),
					Latency: 5 * time.Millisecond,
				}
			}
			return Result{
				NodeID:  node.ID,
				Healthy: false,
				Checked: now(),
				Err:     err,
				Latency: timeout,
			}
		}
		return Result{
			NodeID:  node.ID,
			Healthy: true,
			Checked: now(),
			Latency: latency,
		}
	}

	start := now()
	address := net.JoinHostPort(node.Address, fmt.Sprintf("%d", node.Port))
	conn, err := (&net.Dialer{Timeout: timeout}).DialContext(ctx, "tcp", address)
	if err != nil {
		return Result{
			NodeID:  node.ID,
			Healthy: false,
			Checked: now(),
			Err:     err,
			Latency: timeout,
		}
	}
	_ = conn.Close()

	return Result{
		NodeID:  node.ID,
		Healthy: true,
		Checked: now(),
		Latency: time.Since(start),
	}
}

func pingICMP(ctx context.Context, host string, timeout time.Duration) (time.Duration, error) {
	timeoutSecs := int(timeout.Seconds())
	if timeoutSecs <= 0 {
		timeoutSecs = 1
	}

	cmd := exec.CommandContext(ctx, "ping", "-c", "1", "-W", fmt.Sprintf("%d", timeoutSecs), host)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("ping failed: %w (output: %q)", err, string(output))
	}

	re := regexp.MustCompile(`(?:round-trip|rtt)\s+min/avg/max(?:/mdev)?\s*=\s*[0-9.]+/([0-9.]+)/[0-9.]+`)
	matches := re.FindStringSubmatch(string(output))
	if len(matches) < 2 {
		return 50 * time.Millisecond, nil
	}

	latencyMs, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return 50 * time.Millisecond, nil
	}

	return time.Duration(latencyMs * float64(time.Millisecond)), nil
}

func isCommandNotFound(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return errors.Is(err, exec.ErrNotFound) ||
		strings.Contains(errStr, "executable file not found") ||
		strings.Contains(errStr, "not found in %PATH%") ||
		strings.Contains(errStr, "not found in $PATH")
}
