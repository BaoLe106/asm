package local

import (
	"fmt"
	"os"
	"path/filepath"
)

type WriteLock struct {
	path string
}

func NewWriteLock(layout Layout) *WriteLock {
	return &WriteLock{path: layout.WriteLockPath()}
}

func (l *WriteLock) Lock() error {
	if err := os.MkdirAll(filepath.Dir(l.path), 0o755); err != nil {
		return err
	}
	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return fmt.Errorf("another asm write operation is in progress")
	}
	_ = f.Close()
	return nil
}

func (l *WriteLock) Unlock() {
	_ = os.Remove(l.path)
}
