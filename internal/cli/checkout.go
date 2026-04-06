package cli

import (
	"fmt"

	"github.com/BaoLe106/asm/internal/app"
	"github.com/spf13/cobra"
)

func newCheckoutCmd(opts *rootOptions) *cobra.Command {
	// var version string
	cmd := &cobra.Command{
		Use:	"checkout <version>",
		Short: 	"Checkout to a specific version",
		Args: 	cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			version := args[0]
			svc := app.NewService(opts.repoRoot)
			if err := svc.CheckoutVersion(version); err != nil {
				return err
			}
			fmt.Printf("Checked out version %q\n", version)
			return nil
		},
	}
	// cmd.Flags().StringVar(&version, "version", "", "Version name")
	// _ = cmd.MarkFlagRequired("version")
	return cmd
}
