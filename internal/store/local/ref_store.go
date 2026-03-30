package local

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BaoLe106/asm/internal/domain"
	"github.com/BaoLe106/asm/internal/util"
)

type RefStore struct {
	layout Layout
}

func NewRefStore(layout Layout) *RefStore {
	return &RefStore{layout: layout}
}

func (r *RefStore) readVersions(path string) (domain.VersionsFile, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return domain.VersionsFile{Versions: []domain.VersionRef{}}, nil
		}
		return domain.VersionsFile{}, err
	}
	var vf domain.VersionsFile
	if err := json.Unmarshal(b, &vf); err != nil {
		return domain.VersionsFile{}, err
	}
	return vf, nil
}

func (r *RefStore) writeVersions(path string, vf domain.VersionsFile) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(vf, "", "  ")
	if err != nil {
		return err
	}
	return util.AtomicWriteFile(path, b, 0o644)
}

func (r *RefStore) AppendAgentVersion(agent string, v domain.VersionRef) error {
	path := r.layout.AgentVersionsPath(agent)
	vf, err := r.readVersions(path)
	if err != nil {
		return err
	}
	vf.Versions = append(vf.Versions, v)
	return r.writeVersions(path, vf)
}

func (r *RefStore) AppendSkillsetVersion(skillset string, v domain.VersionRef) error {
	path := r.layout.SkillsetVersionsPath(skillset)
	vf, err := r.readVersions(path)
	if err != nil {
		return err
	}
	vf.Versions = append(vf.Versions, v)
	return r.writeVersions(path, vf)
}

func (r *RefStore) AgentVersions(agent string) ([]domain.VersionRef, error) {
	vf, err := r.readVersions(r.layout.AgentVersionsPath(agent))
	if err != nil {
		return nil, err
	}
	return vf.Versions, nil
}

func (r *RefStore) SkillsetVersions(skillset string) ([]domain.VersionRef, error) {
	vf, err := r.readVersions(r.layout.SkillsetVersionsPath(skillset))
	if err != nil {
		return nil, err
	}
	return vf.Versions, nil
}

func (r *RefStore) DeleteAgentVersionByRelative(agent string, relative int) error {
	path := r.layout.AgentVersionsPath(agent)
	vf, err := r.readVersions(path)
	if err != nil {
		return err
	}
	idx, err := domain.ResolveByRelativeIndex(len(vf.Versions), relative)
	if err != nil {
		return err
	}
	vf.Versions = append(vf.Versions[:idx], vf.Versions[idx+1:]...)
	return r.writeVersions(path, vf)
}

func (r *RefStore) DeleteSkillsetVersionByRelative(skillset string, relative int) error {
	path := r.layout.SkillsetVersionsPath(skillset)
	vf, err := r.readVersions(path)
	if err != nil {
		return err
	}
	idx, err := domain.ResolveByRelativeIndex(len(vf.Versions), relative)
	if err != nil {
		return err
	}
	vf.Versions = append(vf.Versions[:idx], vf.Versions[idx+1:]...)
	return r.writeVersions(path, vf)
}

func (r *RefStore) ListAgentsWithRefs() ([]string, error) {
	entries, err := os.ReadDir(r.layout.AgentRefsRoot())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []string{}, nil
		}
		return nil, err
	}
	out := make([]string, 0)
	for _, e := range entries {
		if e.IsDir() {
			out = append(out, e.Name())
		}
	}
	return out, nil
}

func (r *RefStore) ListSkillsets() ([]string, error) {
	entries, err := os.ReadDir(r.layout.SkillsetRefsRoot())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []string{}, nil
		}
		return nil, err
	}
	out := make([]string, 0)
	for _, e := range entries {
		if e.IsDir() {
			out = append(out, e.Name())
		}
	}
	return out, nil
}

func ValidateExactlyOne(agent, skillset string) error {
	if agent == "" && skillset == "" {
		return fmt.Errorf("must specify either --agent or --skillset")
	}
	if agent != "" && skillset != "" {
		return fmt.Errorf("--agent and --skillset are mutually exclusive")
	}
	return nil
}
