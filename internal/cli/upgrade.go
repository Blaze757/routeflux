package cli

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const routefluxLatestInstallScriptURL = "https://github.com/Blaze757/routeflux/releases/latest/download/install.sh"

// upgradeInstallerPathOverride is set by tests to override the installer path.
var upgradeInstallerPathOverride string

func upgradeInstallerDir() string {
	if upgradeInstallerPathOverride != "" {
		return filepath.Dir(upgradeInstallerPathOverride)
	}
	dir, err := os.MkdirTemp("", "routeflux-upgrade-")
	if err != nil {
		return os.TempDir()
	}
	_ = os.Chmod(dir, 0o700)
	return dir
}

type upgradeResult struct {
	Status         string `json:"status"`
	URL            string `json:"url"`
	ScriptPath     string `json:"script_path"`
	DownloadOutput string `json:"download_output,omitempty"`
	InstallOutput  string `json:"install_output,omitempty"`
}

func runUpgrade(cmd *cobra.Command, jsonOutput bool) error {
	ctx := cmd.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	selfUpdatePath := "/usr/libexec/routeflux-self-update"
	if os.Getenv("ROUTEFLUX_FORCE_UPGRADE") == "" {
		if _, err := os.Stat(selfUpdatePath); err == nil {
			external := exec.CommandContext(ctx, selfUpdatePath)
			var combined bytes.Buffer
			if jsonOutput {
				external.Stdout = &combined
				external.Stderr = &combined
			} else {
				external.Stdout = io.MultiWriter(cmd.OutOrStdout(), &combined)
				external.Stderr = io.MultiWriter(cmd.ErrOrStderr(), &combined)
			}
			if err := external.Run(); err != nil {
				return fmt.Errorf("self-update wrapper: %w", err)
			}

			if jsonOutput {
				output := combined.String()
				status := "ok"
				if strings.Contains(output, "ROUTEFLUX_SELF_UPDATE_STATUS=up-to-date") {
					status = "up-to-date"
				} else if strings.Contains(output, "ROUTEFLUX_SELF_UPDATE_STATUS=updated") {
					status = "updated"
				}

				lines := strings.Split(output, "\n")
				var cleanLines []string
				for _, line := range lines {
					if !strings.HasPrefix(line, "ROUTEFLUX_SELF_UPDATE_STATUS=") {
						cleanLines = append(cleanLines, line)
					}
				}
				cleanMsg := strings.TrimSpace(strings.Join(cleanLines, "\n"))

				upgradeDir := upgradeInstallerDir()
				defer os.RemoveAll(upgradeDir)
				routefluxUpgradeInstallerPath := upgradeInstallerPathOverride
				if routefluxUpgradeInstallerPath == "" {
					routefluxUpgradeInstallerPath = filepath.Join(upgradeDir, "install.sh")
				}

				res := upgradeResult{
					Status:        status,
					URL:           routefluxLatestInstallScriptURL,
					ScriptPath:    routefluxUpgradeInstallerPath,
					InstallOutput: cleanMsg,
				}
				return printOutput(cmd, true, res, "")
			}
			return nil
		}
	}

	upgradeDir := upgradeInstallerDir()
	defer os.RemoveAll(upgradeDir)
	routefluxUpgradeInstallerPath := upgradeInstallerPathOverride
	if routefluxUpgradeInstallerPath == "" {
		routefluxUpgradeInstallerPath = filepath.Join(upgradeDir, "install.sh")
	}

	result := upgradeResult{
		Status:     "ok",
		URL:        routefluxLatestInstallScriptURL,
		ScriptPath: routefluxUpgradeInstallerPath,
	}

	downloadOutput, err := runUpgradeCommand(ctx, cmd, jsonOutput, "wget", "-O", routefluxUpgradeInstallerPath, routefluxLatestInstallScriptURL)
	if err != nil {
		return fmt.Errorf("download latest installer: %w", err)
	}
	result.DownloadOutput = strings.TrimSpace(downloadOutput)

	if upgradeInstallerPathOverride == "" {
		// Download and verify SHA256 checksum
		checksumURL := routefluxLatestInstallScriptURL + ".sha256"
		checksumPath := routefluxUpgradeInstallerPath + ".sha256"
		if _, err := runUpgradeCommand(ctx, cmd, jsonOutput, "wget", "-q", "-O", checksumPath, checksumURL); err == nil {
			expectedSHA256, err := os.ReadFile(checksumPath)
			if err == nil {
				expected := strings.TrimSpace(string(expectedSHA256))
				// Handle "hash  filename" format
				if parts := strings.Fields(expected); len(parts) >= 2 {
					expected = parts[0]
				}
				if err := verifyScriptSHA256(routefluxUpgradeInstallerPath, expected); err != nil {
					return fmt.Errorf("upgrade script verification failed: %w", err)
				}
			}
		} else {
			if jsonOutput {
				fmt.Fprintf(cmd.ErrOrStderr(), "WARNING: could not download checksum, skipping verification\n")
			}
		}

	}
	installOutput, err := runUpgradeCommand(ctx, cmd, jsonOutput, "sh", routefluxUpgradeInstallerPath)
	if err != nil {
		return fmt.Errorf("run latest installer: %w", err)
	}
	result.InstallOutput = strings.TrimSpace(installOutput)

	if jsonOutput {
		return printOutput(cmd, true, result, "")
	}

	return printOutput(cmd, false, nil, "Upgrade completed using the latest release installer.")
}

func runUpgradeCommand(ctx context.Context, cmd *cobra.Command, jsonOutput bool, name string, args ...string) (string, error) {
	external := exec.CommandContext(ctx, name, args...)

	var combined bytes.Buffer
	if jsonOutput {
		external.Stdout = &combined
		external.Stderr = &combined
	} else {
		external.Stdout = io.MultiWriter(cmd.OutOrStdout(), &combined)
		external.Stderr = io.MultiWriter(cmd.ErrOrStderr(), &combined)
	}

	if err := external.Run(); err != nil {
		return combined.String(), fmt.Errorf("%s %s: %w", name, strings.Join(args, " "), err)
	}

	return combined.String(), nil
}

func verifyScriptSHA256(path, expectedSHA256 string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read script: %w", err)
	}
	h := sha256.Sum256(data)
	actual := hex.EncodeToString(h[:])
	if actual != expectedSHA256 {
		return fmt.Errorf("checksum mismatch: expected %s, got %s", expectedSHA256, actual)
	}
	return nil
}
