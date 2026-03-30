package cli

import (
	"fmt"

	"github.com/BaoLe106/asm/internal/app"
	"github.com/spf13/cobra"
)

func newUpdateCmd(opts *rootOptions) *cobra.Command {
	var skillset string
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Create a new version from current state",
		RunE: func(cmd *cobra.Command, args []string) error {
			svc := app.NewService(opts.repoRoot)
			if !cmd.Flags().Changed("skillset") {
				if err := svc.UpdateAgent(); err != nil {
					return err
				}
				fmt.Println("Created new agent version")
				return nil
			}
			useCurrent := skillset == "__CURRENT__"
			name := skillset
			if useCurrent {
				name = ""
			}
			if err := svc.UpdateSkillset(name, useCurrent); err != nil {
				return err
			}
			if useCurrent {
				fmt.Println("Updated current skillset with new version")
			} else {
				fmt.Printf("Updated skillset %q with new version\n", name)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&skillset, "skillset", "", "Skillset name; use --skillset with no value to target current skillset")
	_ = cmd.Flags().Lookup("skillset").NoOptDefVal
	cmd.Flags().Lookup("skillset").NoOptDefVal = "__CURRENT__"
	return cmd
}
