package cli

import (
	"fmt"
	"time"

	"github.com/BaoLe106/asm/internal/app"
	"github.com/spf13/cobra"
)

func newListCmd(opts *rootOptions) *cobra.Command {
	var agent string
	var skillset string
	var listSkill bool
	var listVersion bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List agents, skills, skillsets, or versions",
		RunE: func(cmd *cobra.Command, args []string) error {
			svc := app.NewService(opts.repoRoot)
			switch {
			case listSkill:
				if agent == "" {
					return fmt.Errorf("--skill requires --agent <agent_name>")
				}
				items, err := svc.ListSkillsByAgent(agent)
				if err != nil {
					return err
				}
				for _, it := range items {
					fmt.Println(it)
				}
				return nil
			case listVersion:
				if agent != "" && skillset != "" {
					return fmt.Errorf("--agent and --skillset are mutually exclusive")
				}
				if skillset != "" {
					if skillset == "__ALL__" {
						return fmt.Errorf("--version with --skillset requires a skillset name")
					}
					versions, err := svc.ListSkillsetVersions(skillset)
					if err != nil {
						return err
					}
					for i := len(versions) - 1; i >= 0; i-- {
						relative := len(versions) - 1 - i
						v := versions[i]
						fmt.Printf("version=%d snapshot=%s created=%s note=%s\n", relative, v.SnapshotID, time.Unix(v.CreatedAt, 0).Format(time.RFC3339), v.Note)
					}
					return nil
				}
				if agent == "" {
					return fmt.Errorf("--version requires either --agent <name> or --skillset <name>")
				}
				versions, err := svc.ListAgentVersions(agent)
				if err != nil {
					return err
				}
				for i := len(versions) - 1; i >= 0; i-- {
					relative := len(versions) - 1 - i
					v := versions[i]
					fmt.Printf("version=%d snapshot=%s created=%s note=%s\n", relative, v.SnapshotID, time.Unix(v.CreatedAt, 0).Format(time.RFC3339), v.Note)
				}
				return nil
			default:
				if cmd.Flags().Changed("skillset") {
					if agent != "" {
						return fmt.Errorf("--agent and --skillset are mutually exclusive")
					}
					items, err := svc.ListSkillsets()
					if err != nil {
						return err
					}
					for _, it := range items {
						fmt.Println(it)
					}
					return nil
				}
				if agent != "" {
					fmt.Println(agent)
					return nil
				}
				agents, err := svc.ListAgents()
				if err != nil {
					return err
				}
				for _, a := range agents {
					fmt.Println(a)
				}
				return nil
			}
		},
	}

	cmd.Flags().StringVar(&agent, "agent", "", "Agent name")
	cmd.Flags().BoolVar(&listSkill, "skill", false, "List skills for the given agent")
	cmd.Flags().StringVar(&skillset, "skillset", "", "List skillsets, or provide a skillset name with --version")
	cmd.Flags().Lookup("skillset").NoOptDefVal = "__ALL__"
	cmd.Flags().BoolVar(&listVersion, "version", false, "List versions for the given target")
	return cmd
}
