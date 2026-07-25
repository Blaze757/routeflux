package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func newMoveCmd(opts *rootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "move <subscription-id> <up|down>",
		Short: "Move a subscription up or down in the display order",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			direction := args[1]

			if err := opts.service.MoveSubscription(context.Background(), id, direction); err != nil {
				return err
			}

			return printOutput(
				cmd,
				opts.jsonOutput,
				map[string]string{"moved": id, "direction": direction},
				fmt.Sprintf("Moved subscription %s %s", id, direction),
			)
		},
	}
}
