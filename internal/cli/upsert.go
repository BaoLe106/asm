package cli

import (
	"fmt"

	"github.com/BaoLe106/asm/internal/app"
	"github.com/spf13/cobra"
)

func newUpsertCmd(opts *rootOptions) *cobra.Command {
	// var version string
	cmd := &cobra.Command{
		Use:   	"upsert <version>",
		Short: 	"Upsert a version with current agent/skill changes",
		Args: 	cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			version := args[0]
			svc := app.NewService(opts.repoRoot)
			if err := svc.UpsertVersion(version); err != nil {
				return err
			}
			if version == "" {
				fmt.Println("Upserted current version")
			} else {
				fmt.Printf("Upserted version %q\n", version)
			}
			return nil
		},
	}
	// cmd.Flags().StringVar(&version, "version", "", "Version name (if omitted, upserts current version)")
	return cmd
}
