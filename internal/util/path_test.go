package util

import "testing"

func TestNormalizeSkillList(t *testing.T) {
	in := "crud, review,crud,  ,lint"
	got := NormalizeSkillList(in)
	if len(got) != 3 {
		t.Fatalf("expected 3 skills, got %d", len(got))
	}
	if got[0] != "crud" || got[1] != "review" || got[2] != "lint" {
		t.Fatalf("unexpected parse result: %#v", got)
	}
}
