package app

import "fmt"

func (s *Service) DeleteVersion(agent string, skillset string, version int) error {
	if version < 0 {
		return fmt.Errorf("--version must be >= 0")
	}
	if agent == "" && skillset == "" {
		return fmt.Errorf("must provide --agent or --skillset")
	}
	if agent != "" && skillset != "" {
		return fmt.Errorf("--agent and --skillset are mutually exclusive")
	}
	if agent != "" {
		return s.refs.DeleteAgentVersionByRelative(agent, version)
	}
	return s.refs.DeleteSkillsetVersionByRelative(skillset, version)
}

func (s *Service) CurrentAgentName() (string, error) {
	st, err := s.state.Get()
	if err != nil {
		return "", err
	}
	if st.AgentName != "" {
		return st.AgentName, nil
	}
	agents, err := s.DiscoverAgentFolders()
	if err != nil {
		return "", err
	}
	if len(agents) == 1 {
		return agents[0], nil
	}
	if len(agents) == 0 {
		return "", fmt.Errorf("no agent found in workspace")
	}
	return "", fmt.Errorf("multiple agents found, apply an agent first")
}
