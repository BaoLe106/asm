package app

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/BaoLe106/asm/internal/domain"
	"github.com/BaoLe106/asm/internal/store/local"
)

type SkillSelector func(options []string) ([]string, error)

type Service struct {
	RepoRoot string
	layout   local.Layout
	objects  *local.ObjectStore
	snaps    *local.SnapshotStore
	refs     *local.RefStore
	state    *local.StateStore
	lock     *local.WriteLock
}

func NewService(repoRoot string) *Service {
	layout := local.NewLayout(repoRoot)
	objects := local.NewObjectStore(layout)
	return &Service{
		RepoRoot: repoRoot,
		layout:   layout,
		objects:  objects,
		snaps:    local.NewSnapshotStore(layout, objects),
		refs:     local.NewRefStore(layout),
		state:    local.NewStateStore(layout),
		lock:     local.NewWriteLock(layout),
	}
}

func (s *Service) DiscoverAgentFolders() ([]string, error) {
	entries, err := os.ReadDir(s.RepoRoot)
	if err != nil {
		return nil, err
	}
	out := make([]string, 0)
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasPrefix(name, ".") {
			continue
		}
		if name == ".asm" || name == ".git" {
			continue
		}
		agent := strings.TrimPrefix(name, ".")
		if agent != "" {
			out = append(out, agent)
		}
	}
	sort.Strings(out)
	return out, nil
}

func (s *Service) AgentDirExists(agent string) bool {
	info, err := os.Stat(filepath.Join(s.RepoRoot, "."+agent))
	return err == nil && info.IsDir()
}

func (s *Service) EnsureAgentSnapshot(agent string) error {
	versions, err := s.refs.AgentVersions(agent)
	if err != nil {
		return err
	}
	if len(versions) > 0 {
		return nil
	}
	if !s.AgentDirExists(agent) {
		return fmt.Errorf("agent %q does not exist", agent)
	}
	snap, err := s.snaps.CreateFromAgentDir(agent, domain.SourceAgent, agent, "auto-initialized")
	if err != nil {
		return err
	}
	return s.refs.AppendAgentVersion(agent, domain.VersionRef{SnapshotID: snap.ID, CreatedAt: snap.CreatedAt, Note: "auto-initialized"})
}

func (s *Service) resolveAgentSnapshot(agent string, version int) (domain.Snapshot, string, error) {
	if err := s.EnsureAgentSnapshot(agent); err != nil {
		return domain.Snapshot{}, "", err
	}
	versions, err := s.refs.AgentVersions(agent)
	if err != nil {
		return domain.Snapshot{}, "", err
	}
	idx, err := domain.ResolveByRelativeIndex(len(versions), version)
	if err != nil {
		return domain.Snapshot{}, "", err
	}
	vr := versions[idx]
	snap, err := s.snaps.Get(vr.SnapshotID)
	if err != nil {
		return domain.Snapshot{}, "", err
	}
	return snap, vr.SnapshotID, nil
}

func (s *Service) resolveSkillsetSnapshot(skillset string, version int) (domain.Snapshot, string, error) {
	versions, err := s.refs.SkillsetVersions(skillset)
	if err != nil {
		return domain.Snapshot{}, "", err
	}
	idx, err := domain.ResolveByRelativeIndex(len(versions), version)
	if err != nil {
		return domain.Snapshot{}, "", fmt.Errorf("skillset %q not found or invalid version: %w", skillset, err)
	}
	vr := versions[idx]
	snap, err := s.snaps.Get(vr.SnapshotID)
	if err != nil {
		return domain.Snapshot{}, "", err
	}
	return snap, vr.SnapshotID, nil
}

func nowNote(prefix string) string {
	return fmt.Sprintf("%s at %s", prefix, time.Now().Format(time.RFC3339))
}
