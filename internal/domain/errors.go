package domain

import "fmt"

var (
	ErrNotFound          = fmt.Errorf("not found")
	ErrInvalidVersion    = fmt.Errorf("invalid version")
	ErrNoCurrentAgent    = fmt.Errorf("no current agent")
	ErrNoCurrentSkillset = fmt.Errorf("no current skillset")
)
