package domain

type SourceKind string

const (
	SourceAgent    SourceKind = "agent"
	SourceSkillset SourceKind = "skillset"
)

type SnapshotEntry struct {
	Path string `json:"path"`
	Mode uint32 `json:"mode"`
	Hash string `json:"hash"`
	Size int64  `json:"size"`
}

type Snapshot struct {
	ID         string            `json:"id"`
	SourceKind SourceKind        `json:"source_kind"`
	SourceName string            `json:"source_name"`
	AgentName  string            `json:"agent_name"`
	CreatedAt  int64             `json:"created_at"`
	Entries    []SnapshotEntry   `json:"entries"`
	ExtraMeta  map[string]string `json:"extra_meta,omitempty"`
}

type VersionRef struct {
	SnapshotID string `json:"snapshot_id"`
	CreatedAt  int64  `json:"created_at"`
	Note       string `json:"note,omitempty"`
}

type VersionsFile struct {
	Versions []VersionRef `json:"versions"`
}

type CurrentState struct {
	AgentName     string `json:"agent_name,omitempty"`
	AgentSnapshot string `json:"agent_snapshot,omitempty"`
	SkillsetName  string `json:"skillset_name,omitempty"`
	SkillsetSnap  string `json:"skillset_snapshot,omitempty"`
}
