package local

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/BaoLe106/asm/internal/domain"
	"github.com/BaoLe106/asm/internal/util"
)

type StateStore struct {
	layout Layout
}

func NewStateStore(layout Layout) *StateStore {
	return &StateStore{layout: layout}
}

func (s *StateStore) Get() (domain.CurrentState, error) {
	b, err := os.ReadFile(s.layout.CurrentStatePath())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return domain.CurrentState{}, nil
		}
		return domain.CurrentState{}, err
	}
	var st domain.CurrentState
	if err := json.Unmarshal(b, &st); err != nil {
		return domain.CurrentState{}, err
	}
	return st, nil
}

func (s *StateStore) Save(st domain.CurrentState) error {
	if err := os.MkdirAll(s.layout.StateRoot(), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		return err
	}
	return util.AtomicWriteFile(s.layout.CurrentStatePath(), b, 0o644)
}
