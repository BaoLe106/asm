package cli

import (
	"fmt"

	"github.com/BaoLe106/asm/internal/app"
	"github.com/spf13/cobra"
)

func newStatusCmd(opts *rootOptions) *cobra.Command {
	var agent string
	var skillset string
	// var skillRaw string
	// var version int

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show current agent and skillset status",
		RunE: func(cmd *cobra.Command, args []string) error {
			svc := app.NewService(opts.repoRoot)
			switch {
			case agent == "__ALL__":
				agent, err := svc.GetCurrentAgent()
				if err != nil {
					return err
				}
				fmt.Printf("Current agent: %s\n", agent)
				return nil
			case skillset == "__ALL__":
				skillset, err := svc.GetCurrentSkillset()
				if err != nil {
					return err
				}
				fmt.Printf("Current skillset: %s\n", skillset)
				return nil
			
			default:
				st, err := svc.GetCurrentState()
				if err != nil {
					return err
				}
				fmt.Println(st)
				return nil
			}
		},
	}
	cmd.Flags().StringVar(&agent, "agent", "", "...")
	cmd.Flags().Lookup("agent").NoOptDefVal = "__ALL__"
	cmd.Flags().StringVar(&skillset, "skillset", "", "...")
	cmd.Flags().Lookup("skillset").NoOptDefVal = "__ALL__"
	// cmd.Flags().Lookup("agent").NoOptDefVal = "__ALL__"

	// cmd.Flags().StringVar(&agent, "agent", "", "Show status for the given agent")
	// cmd.Flags().StringVar(&skillset, "skillset", "", "Show status for the given skillset")
	// cmd.Flags().StringVar(&skillRaw, "skill", "", "Show status for the given skill (format: <agent>:<skill>)")
	// cmd.Flags().IntVar(&version, "version", -1, "Show status for the given version (format: <agent|skillset>:<version>)")
	return cmd

}