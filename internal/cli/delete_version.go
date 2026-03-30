package cli

import (
	"fmt"

	"github.com/BaoLe106/asm/internal/app"
	"github.com/spf13/cobra"
)

func newDeleteVersionCmd(opts *rootOptions) *cobra.Command {
	var agent string
	var skillset string
	var version int
	cmd := &cobra.Command{
		Use:   "delete-version",
		Short: "Delete a version by relative index (0=latest)",
		RunE: func(cmd *cobra.Command, args []string) error {
			svc := app.NewService(opts.repoRoot)
			if err := svc.DeleteVersion(agent, skillset, version); err != nil {
				return err
			}
			if agent != "" {
				fmt.Printf("Deleted agent %q version %d\n", agent, version)
			} else {
				fmt.Printf("Deleted skillset %q version %d\n", skillset, version)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&agent, "agent", "", "Agent name")
	cmd.Flags().StringVar(&skillset, "skillset", "", "Skillset name")
	cmd.Flags().IntVar(&version, "version", -1, "Version index to delete (0=latest)")
	_ = cmd.MarkFlagRequired("version")
	return cmd
}
