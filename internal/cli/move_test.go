package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/Alaxay8/routeflux/internal/app"
	"github.com/Alaxay8/routeflux/internal/domain"
)

func TestMoveCommandReordersSubscriptions(t *testing.T) {
	t.Parallel()

	store := &cliMemoryStore{
		subs: []domain.Subscription{
			{ID: "sub-1", DisplayName: "Alpha"},
			{ID: "sub-2", DisplayName: "Beta"},
		},
		settings: domain.DefaultSettings(),
		state:    domain.DefaultRuntimeState(),
	}

	cmd := newMoveCmd(&rootOptions{service: app.NewService(app.Dependencies{Store: store})})
	var stdout bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(new(bytes.Buffer))
	cmd.SetArgs([]string{"sub-2", "up"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute move: %v", err)
	}

	if got := stdout.String(); !strings.Contains(got, "Moved subscription sub-2 up") {
		t.Fatalf("unexpected output: %q", got)
	}
	if len(store.subs) != 2 || store.subs[0].ID != "sub-2" || store.subs[1].ID != "sub-1" {
		t.Fatalf("unexpected order after move command: %+v", store.subs)
	}
}
