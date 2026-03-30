package cli

import (
	"fmt"

	"github.com/BaoLe106/asm/internal/app"
	"github.com/spf13/cobra"
)

func newPublishCmd(opts *rootOptions) *cobra.Command {
	var skillset string
	cmd := &cobra.Command{
		Use:   "publish",
		Short: "Publish current agent state as a skillset version",
		RunE: func(cmd *cobra.Command, args []string) error {
			if skillset == "" {
				return fmt.Errorf("--skillset is required")
			}
			svc := app.NewService(opts.repoRoot)
			if err := svc.PublishSkillset(skillset); err != nil {
				return err
			}
			fmt.Printf("Published skillset %q\n", skillset)
			return nil
		},
	}
	cmd.Flags().StringVar(&skillset, "skillset", "", "Skillset name")
	return cmd
}
