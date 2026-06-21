package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestRestartCommandClearsCaches(t *testing.T) {
	t.Parallel()

	tmpDir, err := os.MkdirTemp("", "routeflux-test-restart")
	if err != nil {
		t.Fatalf("mkdir temp: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(tmpDir) })

	// Create a dummy lock file to verify it gets deleted
	lockPath := filepath.Join(tmpDir, ".routeflux.lock")
	if err := os.WriteFile(lockPath, []byte("lock"), 0600); err != nil {
		t.Fatalf("write lock file: %v", err)
	}

	opts := &rootOptions{
		rootDir: tmpDir,
	}

	cmd := newRestartCmd(opts)
	var stdout bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(new(bytes.Buffer))

	// Note: on testing environments, /etc/init.d/routeflux doesn't exist,
	// so it will fall back to the error output but still clear the caches/locks.
	_ = cmd.Execute()

	if _, err := os.Stat(lockPath); !os.IsNotExist(err) {
		t.Fatalf("expected lock file to be deleted, stat error: %v", err)
	}
}
