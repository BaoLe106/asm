package cli

import (
	"os"

	"github.com/spf13/cobra"
)

type rootOptions struct {
	repoRoot string
}

func NewRootCmd() *cobra.Command {
	opts := &rootOptions{}
	cmd := &cobra.Command{
		Use:           "asm",
		Short:         "AI Agent Skills Manager CLI",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	wd, _ := os.Getwd()
	cmd.PersistentFlags().StringVar(&opts.repoRoot, "repo", wd, "Repository root")
	cmd.AddCommand(newCheckoutCmd(opts))
	cmd.AddCommand(newListCmd(opts))
	cmd.AddCommand(newStatusCmd(opts))
	cmd.AddCommand(newUpsertCmd(opts))
	cmd.AddCommand(newDeleteVersionCmd(opts))
	return cmd
}
