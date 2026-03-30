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
	cmd.AddCommand(newApplyCmd(opts))
	cmd.AddCommand(newListCmd(opts))
	cmd.AddCommand(newPublishCmd(opts))
	cmd.AddCommand(newUpdateCmd(opts))
	cmd.AddCommand(newDeleteVersionCmd(opts))
	return cmd
}
