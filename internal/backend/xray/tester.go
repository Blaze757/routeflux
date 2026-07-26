package xray

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const defaultXrayBinaryPath = "/usr/bin/xray"

// ConfigTester validates rendered Xray configs before they are applied.
type ConfigTester interface {
	Test(ctx context.Context, configPath string) error
}

// CommandTester validates configs by shelling out to `xray -test`.
type CommandTester struct {
	BinaryPath string
}

// NewCommandTester returns a config tester that uses the configured Xray binary.
func NewCommandTester() CommandTester {
	return CommandTester{BinaryPath: xrayBinaryPath()}
}

// BinaryPath returns the configured Xray binary path.
func BinaryPath() string {
	return xrayBinaryPath()
}

// Test validates the provided config path with `xray -test -config`.
func (t CommandTester) Test(ctx context.Context, configPath string) error {
	if strings.TrimSpace(configPath) == "" {
		return fmt.Errorf("xray test config path is required")
	}

	binaryPath := strings.TrimSpace(t.BinaryPath)
	if binaryPath == "" {
		binaryPath = defaultXrayBinaryPath
	}
	if err := validateXrayBinaryPath(binaryPath); err != nil {
		return err
	}

	cmd := exec.CommandContext(ctx, binaryPath, "-test", "-config", configPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("run %s -test -config %s: %w: %s", binaryPath, configPath, err, strings.TrimSpace(string(output)))
	}

	return nil
}

func xrayBinaryPath() string {
	if path := os.Getenv("ROUTEFLUX_XRAY_BINARY"); path != "" {
		return path
	}
	return defaultXrayBinaryPath
}

var allowedXrayBinaryPaths = []string{
	"/usr/bin/xray",
	"/usr/local/bin/xray",
	"/opt/xray/xray",
}

func validateXrayBinaryPath(path string) error {
	if os.Getenv("ROUTEFLUX_INSECURE") != "" {
		return nil
	}
	clean := filepath.Clean(path)
	for _, a := range allowedXrayBinaryPaths {
		if clean == a {
			return nil
		}
	}
	return fmt.Errorf("xray binary path not in allowlist: %s", path)
}
