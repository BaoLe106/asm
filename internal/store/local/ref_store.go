package local

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/BaoLe106/asm/internal/domain"
	"github.com/BaoLe106/asm/internal/util"
)

type RefStore struct {
	layout Layout
}

func NewRefStore(layout Layout) *RefStore {
	return &RefStore{layout: layout}
}

func (r *RefStore) readVersion(path string) (domain.VersionManifest, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return domain.VersionManifest{}, domain.ErrNotFound
		}
		return domain.VersionManifest{}, err
	}
	var vm domain.VersionManifest
	if err := json.Unmarshal(b, &vm); err != nil {
		return domain.VersionManifest{}, err
	}
	return vm, nil
}

func (r *RefStore) writeVersion(path string, vm domain.VersionManifest) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(vm, "", "  ")
	if err != nil {
		return err
	}
	return util.AtomicWriteFile(path, b, 0o644)
}

func (r *RefStore) SaveVersion(name string, vm domain.VersionManifest) error {
	if name == "" {
		return fmt.Errorf("version name is required")
	}
	if vm.Name == "" {
		vm.Name = name
	}
	if vm.Agents == nil {
		vm.Agents = map[string]string{}
	}
	return r.writeVersion(r.layout.VersionPath(name), vm)
}

func (r *RefStore) GetVersion(name string) (domain.VersionManifest, error) {
	if name == "" {
		return domain.VersionManifest{}, fmt.Errorf("version name is required")
	}
	vm, err := r.readVersion(r.layout.VersionPath(name))
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.VersionManifest{}, fmt.Errorf("version %q not found", name)
		}
		return domain.VersionManifest{}, err
	}
	return vm, nil
}

func (r *RefStore) DeleteVersion(name string) error {
	if name == "" {
		return fmt.Errorf("version name is required")
	}
	path := r.layout.VersionPath(name)
	if err := os.Remove(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("version %q not found", name)
		}
		return err
	}
	_ = os.Remove(filepath.Dir(path))
	return nil
}

func (r *RefStore) ListVersions() ([]string, error) {
	entries, err := os.ReadDir(r.layout.VersionRefsRoot())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []string{}, nil
		}
		return nil, err
	}
	out := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			out = append(out, e.Name())
		}
	}
	sort.Strings(out)
	return out, nil
}
