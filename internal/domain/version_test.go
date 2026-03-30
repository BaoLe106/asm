package domain

import "testing"

func TestResolveByRelativeIndex(t *testing.T) {
	tests := []struct {
		total int
		idx   int
		want  int
		ok    bool
	}{
		{total: 3, idx: 0, want: 2, ok: true},
		{total: 3, idx: 1, want: 1, ok: true},
		{total: 3, idx: 2, want: 0, ok: true},
		{total: 3, idx: 3, want: 0, ok: false},
		{total: 0, idx: 0, want: 0, ok: false},
		{total: 2, idx: -1, want: 0, ok: false},
	}

	for _, tc := range tests {
		got, err := ResolveByRelativeIndex(tc.total, tc.idx)
		if tc.ok && err != nil {
			t.Fatalf("unexpected error for total=%d idx=%d: %v", tc.total, tc.idx, err)
		}
		if !tc.ok && err == nil {
			t.Fatalf("expected error for total=%d idx=%d", tc.total, tc.idx)
		}
		if tc.ok && got != tc.want {
			t.Fatalf("want=%d got=%d", tc.want, got)
		}
	}
}
