package app

import (
	"fmt"

	"github.com/BaoLe106/asm/internal/domain"
)

func (s *Service) ApplyAgent(agent string, skills []string, version int, selector SkillSelector) error {
	if agent == "" {
		return fmt.Errorf("--agent is required")
	}
	if !s.AgentDirExists(agent) {
		return fmt.Errorf("agent %q does not exist", agent)
	}

	snap, snapID, err := s.resolveAgentSnapshot(agent, version)
	if err != nil {
		return err
	}
	available := s.snaps.ListSkills(snap)
	selected := make([]string, 0)

	if len(skills) == 0 {
		if selector != nil {
			chosen, err := selector(available)
			if err != nil {
				return err
			}
			selected = chosen
		}
	} else {
		selected = skills
	}

	// Empty selection means apply all skills.
	selectedSet := map[string]bool{}
	if len(selected) > 0 {
		valid := map[string]bool{}
		for _, sk := range available {
			valid[sk] = true
		}
		for _, sk := range selected {
			if !valid[sk] {
				return fmt.Errorf("skill %q does not exist for agent %q", sk, agent)
			}
			selectedSet[sk] = true
		}
	}

	if err := s.lock.Lock(); err != nil {
		return err
	}
	defer s.lock.Unlock()

	if err := s.snaps.RestoreToAgentDir(snap, selectedSet); err != nil {
		return err
	}

	st, _ := s.state.Get()
	st.AgentName = agent
	st.AgentSnapshot = snapID
	if err := s.state.Save(st); err != nil {
		return err
	}
	return nil
}

func (s *Service) ApplySkillset(skillset string, version int) error {
	if skillset == "" {
		return fmt.Errorf("--skillset is required")
	}
	snap, snapID, err := s.resolveSkillsetSnapshot(skillset, version)
	if err != nil {
		return err
	}
	if err := s.lock.Lock(); err != nil {
		return err
	}
	defer s.lock.Unlock()
	if err := s.snaps.RestoreToAgentDir(snap, map[string]bool{}); err != nil {
		return err
	}
	st, _ := s.state.Get()
	st.AgentName = snap.AgentName
	st.AgentSnapshot = snapID
	st.SkillsetName = skillset
	st.SkillsetSnap = snapID
	return s.state.Save(st)
}

func (s *Service) PublishSkillset(skillsetName string) error {
	if skillsetName == "" {
		return fmt.Errorf("--skillset is required")
	}
	agent, err := s.CurrentAgentName()
	if err != nil {
		return err
	}
	snap, err := s.snaps.CreateFromAgentDir(agent, domain.SourceSkillset, skillsetName, nowNote("publish"))
	if err != nil {
		return err
	}
	if err := s.refs.AppendSkillsetVersion(skillsetName, domain.VersionRef{SnapshotID: snap.ID, CreatedAt: snap.CreatedAt, Note: "publish"}); err != nil {
		return err
	}
	st, _ := s.state.Get()
	st.SkillsetName = skillsetName
	st.SkillsetSnap = snap.ID
	return s.state.Save(st)
}
