package local

import "path/filepath"

type Layout struct {
	RepoRoot string
}

func NewLayout(repoRoot string) Layout {
	return Layout{RepoRoot: repoRoot}
}

func (l Layout) AsmRoot() string          { return filepath.Join(l.RepoRoot, ".asm") }
func (l Layout) ObjectsRoot() string      { return filepath.Join(l.AsmRoot(), "objects") }
func (l Layout) SnapshotsRoot() string    { return filepath.Join(l.AsmRoot(), "snapshots") }
func (l Layout) RefsRoot() string         { return filepath.Join(l.AsmRoot(), "refs") }
func (l Layout) AgentRefsRoot() string    { return filepath.Join(l.RefsRoot(), "agents") }
func (l Layout) SkillsetRefsRoot() string { return filepath.Join(l.RefsRoot(), "skillsets") }
func (l Layout) StateRoot() string        { return filepath.Join(l.AsmRoot(), "state") }
func (l Layout) CurrentStatePath() string { return filepath.Join(l.StateRoot(), "current.json") }
func (l Layout) LocksRoot() string        { return filepath.Join(l.AsmRoot(), "locks") }
func (l Layout) WriteLockPath() string    { return filepath.Join(l.LocksRoot(), "write.lock") }
func (l Layout) AgentVersionsPath(name string) string {
	return filepath.Join(l.AgentRefsRoot(), name, "versions.json")
}
func (l Layout) SkillsetVersionsPath(name string) string {
	return filepath.Join(l.SkillsetRefsRoot(), name, "versions.json")
}
func (l Layout) SnapshotPath(id string) string {
	return filepath.Join(l.SnapshotsRoot(), id+".json")
}
