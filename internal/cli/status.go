package cli

import (
	"fmt"

	"github.com/BaoLe106/asm/internal/app"
	"github.com/spf13/cobra"
)

func newStatusCmd(opts *rootOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show current version status",
		RunE: func(cmd *cobra.Command, args []string) error {
			svc := app.NewService(opts.repoRoot)
			version, err := svc.GetCurrentVersion()
			if err != nil {
				return err
			}
			fmt.Printf("Current version: %s\n", version)
			return nil
		},
	}
	return cmd
}
