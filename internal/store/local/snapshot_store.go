package local

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/BaoLe106/asm/internal/domain"
	"github.com/BaoLe106/asm/internal/util"
)

type SnapshotStore struct {
	layout  Layout
	objects *ObjectStore
}

func NewSnapshotStore(layout Layout, objects *ObjectStore) *SnapshotStore {
	return &SnapshotStore{layout: layout, objects: objects}
}

func (s *SnapshotStore) CreateFromAgentDir(agentName string, sourceKind domain.SourceKind, sourceName string, note string) (domain.Snapshot, error) {
	agentDir := filepath.Join(s.layout.RepoRoot, "."+agentName)
	info, err := os.Stat(agentDir)
	if err != nil || !info.IsDir() {
		return domain.Snapshot{}, fmt.Errorf("agent directory .%s not found", agentName)
	}

	entries := make([]domain.SnapshotEntry, 0)
	err = filepath.WalkDir(agentDir, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(agentDir, path)
		if err != nil {
			return err
		}
		rel = util.ToSlash(rel)
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		h, err := s.objects.Put(data)
		if err != nil {
			return err
		}
		fi, err := os.Stat(path)
		if err != nil {
			return err
		}
		entries = append(entries, domain.SnapshotEntry{
			Path: rel,
			Mode: uint32(fi.Mode().Perm()),
			Hash: h,
			Size: fi.Size(),
		})
		return nil
	})
	if err != nil {
		return domain.Snapshot{}, err
	}

	sort.Slice(entries, func(i, j int) bool { return entries[i].Path < entries[j].Path })
	now := time.Now().Unix()
	id := fmt.Sprintf("snp_%d_%s", now, util.SHA256([]byte(fmt.Sprintf("%s:%s:%d", sourceKind, sourceName, now)))[:8])
	snap := domain.Snapshot{
		ID:         id,
		SourceKind: sourceKind,
		SourceName: sourceName,
		AgentName:  agentName,
		CreatedAt:  now,
		Entries:    entries,
		ExtraMeta: map[string]string{
			"note": note,
		},
	}
	if err := s.Save(snap); err != nil {
		return domain.Snapshot{}, err
	}
	return snap, nil
}

func (s *SnapshotStore) Save(snap domain.Snapshot) error {
	if err := os.MkdirAll(s.layout.SnapshotsRoot(), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return err
	}
	return util.AtomicWriteFile(s.layout.SnapshotPath(snap.ID), b, 0o644)
}

func (s *SnapshotStore) Get(id string) (domain.Snapshot, error) {
	b, err := os.ReadFile(s.layout.SnapshotPath(id))
	if err != nil {
		return domain.Snapshot{}, err
	}
	var snap domain.Snapshot
	if err := json.Unmarshal(b, &snap); err != nil {
		return domain.Snapshot{}, err
	}
	return snap, nil
}

func (s *SnapshotStore) ListSkills(snap domain.Snapshot) []string {
	set := map[string]struct{}{}
	for _, e := range snap.Entries {
		if strings.HasPrefix(e.Path, "skills/") {
			parts := strings.Split(e.Path, "/")
			if len(parts) >= 2 && parts[1] != "" {
				set[parts[1]] = struct{}{}
			}
		}
	}
	out := make([]string, 0, len(set))
	for k := range set {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func (s *SnapshotStore) RestoreToAgentDir(snap domain.Snapshot, selectedSkills map[string]bool) error {
	agentDir := filepath.Join(s.layout.RepoRoot, "."+snap.AgentName)
	if err := util.RemoveAndCreateDir(agentDir); err != nil {
		return err
	}
	for _, e := range snap.Entries {
		if strings.HasPrefix(e.Path, "skills/") && len(selectedSkills) > 0 {
			parts := strings.Split(e.Path, "/")
			if len(parts) >= 2 {
				if !selectedSkills[parts[1]] {
					continue
				}
			}
		}
		data, err := s.objects.Get(e.Hash)
		if err != nil {
			return err
		}
		outPath := filepath.Join(agentDir, filepath.FromSlash(e.Path))
		if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(outPath, data, os.FileMode(e.Mode)); err != nil {
			return err
		}
	}
	return nil
}
