package app

import "github.com/BaoLe106/asm/internal/domain"

func (s *Service) GetCurrentState() (domain.CurrentState, error) {
	return s.state.Get()
}

func (s *Service) GetCurrentVersion() (string, error) {
	return s.readCurrentVersionName()
}
