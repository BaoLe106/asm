package app

import (
	"fmt"

	"github.com/BaoLe106/asm/internal/domain"
)

func (s *Service) GetCurrentState() (domain.CurrentState, error) {
	return s.state.Get()
}

func (s *Service) GetCurrentAgent() (string, error) {
	st, err := s.state.Get()
	if err != nil {
		return "", err
	}
	if st.AgentName == "" {
		return "", fmt.Errorf("%w", domain.ErrNoCurrentAgent)
	}
	return st.AgentName, nil
}

func (s *Service) GetCurrentSkillset() (string, error) {
	st, err := s.state.Get()
	if err != nil {
		return "", err
	}
	if st.SkillsetName == "" {
		return "", fmt.Errorf("%w", domain.ErrNoCurrentSkillset)
	}
	return st.SkillsetName, nil
}