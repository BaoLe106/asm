package app

import "fmt"

func (s *Service) CheckoutVersion(version string) error {
	if version == "" {
		return fmt.Errorf("--version is required")
	}
	manifest, err := s.refs.GetVersion(version)
	if err != nil {
		return err
	}
	if err := s.lock.Lock(); err != nil {
		return err
	}
	defer s.lock.Unlock()

	if err := s.removeMissingAgents(manifest.Agents); err != nil {
		return err
	}
	for agent, snapID := range manifest.Agents {
		snap, err := s.snaps.Get(snapID)
		if err != nil {
			return err
		}
		snap.AgentName = agent
		if err := s.snaps.RestoreToAgentDir(snap, map[string]bool{}); err != nil {
			return err
		}
	}
	return s.setCurrentVersionName(version)
}

func (s *Service) CreateVersion(name string) error {
	if name == "" {
		return fmt.Errorf("--version is required")
	}
	manifest, err := s.captureVersionManifest(name, "create version")
	if err != nil {
		return err
	}
	if err := s.refs.SaveVersion(name, manifest); err != nil {
		return err
	}
	return s.setCurrentVersionName(name)
}
