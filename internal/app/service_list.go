package app

import (
	"fmt"
	"sort"
)

func (s *Service) ListVersions() ([]string, error) {
	return s.refs.ListVersions()
}

func (s *Service) ListCurrentVersionAgents() ([]string, error) {
	current, err := s.readCurrentVersionName()
	if err != nil {
		return nil, err
	}
	manifest, err := s.refs.GetVersion(current)
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(manifest.Agents))
	for agent := range manifest.Agents {
		out = append(out, agent)
	}
	sort.Strings(out)
	return out, nil
}

func (s *Service) ListCurrentVersionSkillsByAgent(agent string) ([]string, error) {
	if agent == "" {
		return nil, fmt.Errorf("--agent is required")
	}
	current, err := s.readCurrentVersionName()
	if err != nil {
		return nil, err
	}
	manifest, err := s.refs.GetVersion(current)
	if err != nil {
		return nil, err
	}
	snapID, ok := manifest.Agents[agent]
	if !ok {
		return nil, fmt.Errorf("agent %q not found in current version %q", agent, current)
	}
	snap, err := s.snaps.Get(snapID)
	if err != nil {
		return nil, err
	}
	return s.snaps.ListSkills(snap), nil
}

func (s *Service) ListAllSkillsInCurrentVersion() (map[string][]string, error) {
	current, err := s.readCurrentVersionName()
	if err != nil {
		return nil, err
	}
	manifest, err := s.refs.GetVersion(current)
	if err != nil {
		return nil, err
	}
	out := map[string][]string{}
	for agent, snapID := range manifest.Agents {
		snap, err := s.snaps.Get(snapID)
		if err != nil {
			return nil, err
		}
		out[agent] = s.snaps.ListSkills(snap)
	}
	return out, nil
}
