package local

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BaoLe106/asm/internal/util"
)

type ObjectStore struct {
	layout Layout
}

func NewObjectStore(layout Layout) *ObjectStore {
	return &ObjectStore{layout: layout}
}

func (s *ObjectStore) Put(data []byte) (string, error) {
	h := util.SHA256(data)
	dir := filepath.Join(s.layout.ObjectsRoot(), h[:2])
	path := filepath.Join(dir, h)
	if _, err := os.Stat(path); err == nil {
		return h, nil
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	if err := util.AtomicWriteFile(path, data, 0o644); err != nil {
		return "", err
	}
	return h, nil
}

func (s *ObjectStore) Get(hash string) ([]byte, error) {
	if len(hash) < 2 {
		return nil, fmt.Errorf("invalid hash")
	}
	path := filepath.Join(s.layout.ObjectsRoot(), hash[:2], hash)
	return os.ReadFile(path)
}
