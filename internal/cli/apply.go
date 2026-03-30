package cli

import (
	"fmt"

	"github.com/BaoLe106/asm/internal/app"
	"github.com/BaoLe106/asm/internal/ui"
	"github.com/BaoLe106/asm/internal/util"
	"github.com/spf13/cobra"
)

func newApplyCmd(opts *rootOptions) *cobra.Command {
	var agent string
	var skillset string
	var skillRaw string
	var version int

	cmd := &cobra.Command{
		Use:   "apply",
		Short: "Apply agent or skillset configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			svc := app.NewService(opts.repoRoot)
			if agent != "" && skillset != "" {
				return fmt.Errorf("--agent and --skillset are mutually exclusive")
			}
			if agent == "" && skillset == "" {
				return fmt.Errorf("must provide --agent or --skillset")
			}
			if skillset != "" {
				if skillRaw != "" {
					return fmt.Errorf("--skill can only be used with --agent")
				}
				if err := svc.ApplySkillset(skillset, version); err != nil {
					return err
				}
				fmt.Printf("Applied skillset %q version %d\n", skillset, version)
				return nil
			}
			skills := util.NormalizeSkillList(skillRaw)
			if err := svc.ApplyAgent(agent, skills, version, ui.SelectSkills); err != nil {
				return err
			}
			fmt.Printf("Applied agent %q version %d\n", agent, version)
			return nil
		},
	}

	cmd.Flags().StringVar(&agent, "agent", "", "Agent name")
	cmd.Flags().StringVar(&skillset, "skillset", "", "Skillset name")
	cmd.Flags().StringVar(&skillRaw, "skill", "", "Comma-separated skill names (e.g. crud,review)")
	cmd.Flags().IntVar(&version, "version", 0, "Version index (0=latest, larger is older)")
	return cmd
}
