package ui

import (
	"github.com/AlecAivazis/survey/v2"
)

func SelectSkills(options []string) ([]string, error) {
	if len(options) == 0 {
		return []string{}, nil
	}
	selected := []string{}
	prompt := &survey.MultiSelect{
		Message: "Select skills (Space to toggle, Enter with no selection means all):",
		Options: options,
	}
	err := survey.AskOne(prompt, &selected)
	if err != nil {
		return nil, err
	}
	return selected, nil
}
