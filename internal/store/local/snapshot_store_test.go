package local

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/BaoLe106/asm/internal/domain"
)

func TestCreateFromAgentDirSkipsGitHubWorkflows(t *testing.T) {
	repoRoot := t.TempDir()
	store := NewSnapshotStore(NewLayout(repoRoot), NewObjectStore(NewLayout(repoRoot)))

	if err := os.MkdirAll(filepath.Join(repoRoot, ".github", "workflows"), 0o755); err != nil {
		t.Fatalf("mkdir workflows: %v", err)
	}
	if err := os.WriteFile(filepath.Join(repoRoot, ".github", "workflows", "ci.yml"), []byte("name: ci"), 0o644); err != nil {
		t.Fatalf("write workflow: %v", err)
	}
	if err := os.WriteFile(filepath.Join(repoRoot, ".github", "CODEOWNERS"), []byte("* @team"), 0o644); err != nil {
		t.Fatalf("write CODEOWNERS: %v", err)
	}

	snap, err := store.CreateFromAgentDir("github", domain.SourceAgent, "v1", "test")
	if err != nil {
		t.Fatalf("create snapshot: %v", err)
	}

	for _, entry := range snap.Entries {
		if entry.Path == "workflows/ci.yml" {
			t.Fatalf("workflow file should not be recorded: %s", entry.Path)
		}
	}

	if len(snap.Entries) != 1 || snap.Entries[0].Path != "CODEOWNERS" {
		t.Fatalf("unexpected snapshot entries: %#v", snap.Entries)
	}
}

func TestRestoreToAgentDirPreservesGitHubWorkflows(t *testing.T) {
	repoRoot := t.TempDir()
	layout := NewLayout(repoRoot)
	objects := NewObjectStore(layout)
	store := NewSnapshotStore(layout, objects)

	if err := os.MkdirAll(filepath.Join(repoRoot, ".github", "workflows"), 0o755); err != nil {
		t.Fatalf("mkdir workflows: %v", err)
	}
	if err := os.WriteFile(filepath.Join(repoRoot, ".github", "workflows", "legacy.yml"), []byte("name: legacy"), 0o644); err != nil {
		t.Fatalf("write legacy workflow: %v", err)
	}

	readmeHash, err := objects.Put([]byte("hello"))
	if err != nil {
		t.Fatalf("put readme object: %v", err)
	}
	workflowHash, err := objects.Put([]byte("name: should-not-restore"))
	if err != nil {
		t.Fatalf("put workflow object: %v", err)
	}

	snap := domain.Snapshot{
		AgentName: "github",
		Entries: []domain.SnapshotEntry{
			{Path: "README.md", Mode: 0o644, Hash: readmeHash},
			{Path: "workflows/new.yml", Mode: 0o644, Hash: workflowHash},
		},
	}

	if err := store.RestoreToAgentDir(snap, map[string]bool{}); err != nil {
		t.Fatalf("restore snapshot: %v", err)
	}

	if _, err := os.Stat(filepath.Join(repoRoot, ".github", "workflows")); err != nil {
		t.Fatalf("workflows dir should exist: %v", err)
	}

	legacyContent, err := os.ReadFile(filepath.Join(repoRoot, ".github", "workflows", "legacy.yml"))
	if err != nil {
		t.Fatalf("legacy workflow should be preserved: %v", err)
	}
	if string(legacyContent) != "name: legacy" {
		t.Fatalf("legacy workflow content changed: %q", string(legacyContent))
	}

	if _, err := os.Stat(filepath.Join(repoRoot, ".github", "workflows", "new.yml")); !os.IsNotExist(err) {
		t.Fatalf("workflow from snapshot should not be restored")
	}

	readmeContent, err := os.ReadFile(filepath.Join(repoRoot, ".github", "README.md"))
	if err != nil {
		t.Fatalf("README should be restored: %v", err)
	}
	if string(readmeContent) != "hello" {
		t.Fatalf("unexpected README content: %q", string(readmeContent))
	}
}
