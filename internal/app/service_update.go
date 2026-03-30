package app

import (
	"fmt"

	"github.com/BaoLe106/asm/internal/domain"
)

func (s *Service) UpdateAgent() error {
	agent, err := s.CurrentAgentName()
	if err != nil {
		return err
	}
	snap, err := s.snaps.CreateFromAgentDir(agent, domain.SourceAgent, agent, nowNote("update"))
	if err != nil {
		return err
	}
	if err := s.refs.AppendAgentVersion(agent, domain.VersionRef{SnapshotID: snap.ID, CreatedAt: snap.CreatedAt, Note: "update"}); err != nil {
		return err
	}
	st, _ := s.state.Get()
	st.AgentName = agent
	st.AgentSnapshot = snap.ID
	return s.state.Save(st)
}

func (s *Service) UpdateSkillset(name string, useCurrent bool) error {
	target := name
	if useCurrent {
		st, err := s.state.Get()
		if err != nil {
			return err
		}
		if st.SkillsetName == "" {
			return fmt.Errorf("%w", domain.ErrNoCurrentSkillset)
		}
		target = st.SkillsetName
	}
	if target == "" {
		return fmt.Errorf("skillset name is required")
	}
	versions, err := s.refs.SkillsetVersions(target)
	if err != nil {
		return err
	}
	if len(versions) == 0 {
		return fmt.Errorf("skillset %q does not exist", target)
	}
	agent, err := s.CurrentAgentName()
	if err != nil {
		return err
	}
	snap, err := s.snaps.CreateFromAgentDir(agent, domain.SourceSkillset, target, nowNote("update skillset"))
	if err != nil {
		return err
	}
	if err := s.refs.AppendSkillsetVersion(target, domain.VersionRef{SnapshotID: snap.ID, CreatedAt: snap.CreatedAt, Note: "update"}); err != nil {
		return err
	}
	st, _ := s.state.Get()
	st.SkillsetName = target
	st.SkillsetSnap = snap.ID
	return s.state.Save(st)
}
