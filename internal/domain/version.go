package domain

import "fmt"

// ResolveByRelativeIndex maps a user index (0=latest, 1=previous) to slice index.
func ResolveByRelativeIndex(total int, idx int) (int, error) {
	if idx < 0 {
		return 0, fmt.Errorf("%w: version must be >= 0", ErrInvalidVersion)
	}
	if total == 0 {
		return 0, fmt.Errorf("%w: no versions available", ErrInvalidVersion)
	}
	absolute := total - 1 - idx
	if absolute < 0 || absolute >= total {
		return 0, fmt.Errorf("%w: version %d is out of range", ErrInvalidVersion, idx)
	}
	return absolute, nil
}
