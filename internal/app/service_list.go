package app

import (
	"fmt"
	"sort"

	"github.com/BaoLe106/asm/internal/domain"
)

func (s *Service) ListAgents() ([]string, error) {
	dirs, err := s.DiscoverAgentFolders()
	if err != nil {
		return nil, err
	}
	return dirs, nil
}

func (s *Service) ListSkillsets() ([]string, error) {
	items, err := s.refs.ListSkillsets()
	if err != nil {
		return nil, err
	}
	sort.Strings(items)
	return items, nil
}

func (s *Service) ListSkillsByAgent(agent string) ([]string, error) {
	if agent == "" {
		return nil, fmt.Errorf("--agent requires a value")
	}
	snap, _, err := s.resolveAgentSnapshot(agent, 0)
	if err != nil {
		return nil, err
	}
	return s.snaps.ListSkills(snap), nil
}

func (s *Service) ListAgentVersions(agent string) ([]domain.VersionRef, error) {
	if err := s.EnsureAgentSnapshot(agent); err != nil {
		return nil, err
	}
	return s.refs.AgentVersions(agent)
}

func (s *Service) ListSkillsetVersions(skillset string) ([]domain.VersionRef, error) {
	return s.refs.SkillsetVersions(skillset)
}
