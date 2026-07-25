package openwrt

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// allowedNFTPaths defines acceptable nft binary locations.
var allowedNFTPaths = []string{
	"/usr/sbin/nft",
	"/sbin/nft",
}

// allowedIPPaths defines acceptable ip binary locations.
var allowedIPPaths = []string{
	"/sbin/ip",
	"/usr/sbin/ip",
	"/bin/ip",
	"/usr/bin/ip",
}

// allowedXrayPaths defines acceptable xray binary locations.
var allowedXrayPaths = []string{
	"/usr/bin/xray",
	"/usr/local/bin/xray",
	"/opt/xray/xray",
}

// allowedDNSMasqPaths defines acceptable dnsmasq binary locations.
var allowedDNSMasqPaths = []string{
	"/usr/sbin/dnsmasq",
	"/sbin/dnsmasq",
}

// allowedDNSMasqServicePaths defines acceptable dnsmasq init.d locations.
var allowedDNSMasqServicePaths = []string{
	"/etc/init.d/dnsmasq",
}

// allowedZapretServicePaths defines acceptable zapret init.d locations.
var allowedZapretServicePaths = []string{
	"/etc/init.d/zapret",
	"/etc/init.d/zapret-openwrt",
}

// ValidateBinaryPath checks that a path is in the allowlist.
func ValidateBinaryPath(path string, allowed []string) error {
	clean := filepath.Clean(path)
	for _, a := range allowed {
		if clean == a {
			return nil
		}
	}
	return fmt.Errorf("path not in allowlist: %s", path)
}

// ValidateNFTPath validates the nft binary path.
func ValidateNFTPath(path string) error {
	return ValidateBinaryPath(path, allowedNFTPaths)
}

// ValidateIPPath validates the ip binary path.
func ValidateIPPath(path string) error {
	return ValidateBinaryPath(path, allowedIPPaths)
}

// ValidateXrayPath validates the xray binary path.
func ValidateXrayPath(path string) error {
	return ValidateBinaryPath(path, allowedXrayPaths)
}

// ValidateDNSMasqPath validates the dnsmasq binary path.
func ValidateDNSMasqPath(path string) error {
	return ValidateBinaryPath(path, allowedDNSMasqPaths)
}

// ValidateDNSMasqServicePath validates the dnsmasq service path.
func ValidateDNSMasqServicePath(path string) error {
	return ValidateBinaryPath(path, allowedDNSMasqServicePaths)
}

// ValidateZapretServicePath validates the zapret service path.
func ValidateZapretServicePath(path string) error {
	return ValidateBinaryPath(path, allowedZapretServicePaths)
}

// ValidateRootDir validates that ROUTEFLUX_ROOT is within allowed directories.
func ValidateRootDir(root string) error {
	clean := filepath.Clean(root)
	allowedBases := []string{
		"/etc/routeflux",
		"/opt/routeflux",
		"/var/lib/routeflux",
	}
	for _, base := range allowedBases {
		if clean == base || strings.HasPrefix(clean, base+"/") {
			return nil
		}
	}
	return fmt.Errorf("root directory not in allowlist: %s", clean)
}

// ValidateHostname checks that a hostname is a valid IPv4, IPv6, or RFC 1123 hostname.
func ValidateHostname(host string) error {
	if ip := net.ParseIP(host); ip != nil {
		return nil
	}

	if len(host) > 253 {
		return fmt.Errorf("hostname too long: %s", host)
	}

	labels := strings.Split(host, ".")
	for _, label := range labels {
		if len(label) == 0 || len(label) > 63 {
			return fmt.Errorf("invalid label length in hostname: %s", host)
		}
		if strings.HasPrefix(label, "-") || strings.HasSuffix(label, "-") {
			return fmt.Errorf("label cannot start or end with hyphen: %s", label)
		}
		for _, r := range label {
			if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-') {
				return fmt.Errorf("invalid character %c in hostname label: %s", r, label)
			}
		}
	}

	reserved := []string{".local", ".internal", ".test", ".localhost", ".onion"}
	lower := strings.ToLower(host)
	for _, tld := range reserved {
		if strings.HasSuffix(lower, tld) || strings.HasSuffix(lower, tld+".") {
			return fmt.Errorf("reserved TLD in hostname: %s", host)
		}
	}

	return nil
}

var validHostnameRegex = regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?)*$`)

// EnsureDir creates a directory with restricted permissions.
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0o700)
}

// PrivateFilePerm is the permission for sensitive config files.
const PrivateFilePerm = 0o600

// PublicFilePerm is the permission for non-sensitive files.
const PublicFilePerm = 0o644
