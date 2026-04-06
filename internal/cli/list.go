package cli

import (
	"fmt"
	"sort"

	"github.com/BaoLe106/asm/internal/app"
	"github.com/spf13/cobra"
)

func newListCmd(opts *rootOptions) *cobra.Command {
	var agent string
	var listVersion bool
	var listAgent bool
	var listSkill bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List versions or current-version agents/skills",
		RunE: func(cmd *cobra.Command, args []string) error {
			svc := app.NewService(opts.repoRoot)

			switch {
			case listVersion:
				versions, err := svc.ListVersions()
				if err != nil {
					return err
				}
				for _, v := range versions {
					fmt.Println(v)
				}
				return nil
			case listAgent:
				agents, err := svc.ListCurrentVersionAgents()
				if err != nil {
					return err
				}
				for _, a := range agents {
					fmt.Println(a)
				}
				return nil
			case listSkill:
				if agent != "" {
					skills, err := svc.ListCurrentVersionSkillsByAgent(agent)
					if err != nil {
						return err
					}
					for _, sk := range skills {
						fmt.Println(sk)
					}
					return nil
				}
				all, err := svc.ListAllSkillsInCurrentVersion()
				if err != nil {
					return err
				}
				agents := make([]string, 0, len(all))
				for a := range all {
					agents = append(agents, a)
				}
				sort.Strings(agents)
				for _, a := range agents {
					for _, sk := range all[a] {
						fmt.Printf("%s:%s\n", a, sk)
					}
				}
				return nil
			default:
				return fmt.Errorf("use one of: --version, --agent, or --skill")
			}
		},
	}

	cmd.Flags().BoolVar(&listVersion, "version", false, "List version names")
	cmd.Flags().BoolVar(&listAgent, "agent", false, "List all agents in current version")
	cmd.Flags().BoolVar(&listSkill, "skill", false, "List skills in current version")
	cmd.Flags().StringVar(&agent, "agent-name", "", "When used with --skill, list skills for one agent")
	return cmd
}
