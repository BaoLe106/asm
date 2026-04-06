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

func (s *Service) discoverAgentFolders() ([]string, error) {
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
		if !strings.HasPrefix(name, ".") || name == ".asm" || name == ".git" {
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

func (s *Service) captureVersionManifest(name string, note string) (domain.VersionManifest, error) {
	agents, err := s.discoverAgentFolders()
	if err != nil {
		return domain.VersionManifest{}, err
	}
	manifest := domain.VersionManifest{
		Name:      name,
		CreatedAt: time.Now().Unix(),
		Note:      note,
		Agents:    map[string]string{},
	}
	for _, agent := range agents {
		snap, err := s.snaps.CreateFromAgentDir(agent, domain.SourceAgent, name, note)
		if err != nil {
			return domain.VersionManifest{}, err
		}
		manifest.Agents[agent] = snap.ID
	}
	return manifest, nil
}

func (s *Service) readCurrentVersionName() (string, error) {
	st, err := s.state.Get()
	if err != nil {
		return "", err
	}
	if st.CurrentVersion == "" {
		return "", fmt.Errorf("no current version")
	}
	return st.CurrentVersion, nil
}

func (s *Service) setCurrentVersionName(version string) error {
	st, err := s.state.Get()
	if err != nil {
		return err
	}
	st.CurrentVersion = version
	return s.state.Save(st)
}

func (s *Service) removeMissingAgents(keep map[string]string) error {
	agents, err := s.discoverAgentFolders()
	if err != nil {
		return err
	}
	for _, agent := range agents {
		if _, ok := keep[agent]; ok {
			continue
		}
		if err := os.RemoveAll(filepath.Join(s.RepoRoot, "."+agent)); err != nil {
			return err
		}
	}
	return nil
}
