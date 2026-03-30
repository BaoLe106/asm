# asm

AI Agent Skills Manager CLI written in Go.

`asm` manages local AI agent setups in a repository using a Git-like hidden store (`.asm`).

## Key Concepts

- Working tree: `.{agent}/skills/{skillName}` plus any additional files under `.{agent}`
- Local metadata/object store: `.asm/`
- Versioning model:
  - `--version 0` = latest
  - `--version 1` = previous
  - larger number = older

## Local Storage Layout

```
.asm/
  state/current.json
  objects/<hash-prefix>/<hash>
  snapshots/<snapshot-id>.json
  refs/agents/<agent>/versions.json
  refs/skillsets/<skillset>/versions.json
  locks/write.lock
```

## Commands

### Apply

```bash
asm apply --agent <agent_name>
asm apply --agent <agent_name> --skill crud_skill,code_review_skill
asm apply --agent <agent_name> --version 0
asm apply --skillset <skillset_name>
asm apply --skillset <skillset_name> --version 0
```

Behavior:

- If `--skill` is omitted, interactive skill picker opens (arrow keys + space).
- Pressing Enter with no selection applies all skills.

### List

```bash
asm list --agent
asm list --agent <agent_name> --skill
asm list --skillset
asm list --agent <agent_name> --version
asm list --skillset <skillset_name> --version
```

### Publish

```bash
asm publish --skillset <skillset_name>
```

Creates a new skillset version from current `.{agent}` folder state.

### Update

```bash
asm update
asm update --skillset
asm update --skillset <skillset_name>
```

- `asm update`: create new version for current agent.
- `asm update --skillset`: update current skillset (from state), error if none.
- `asm update --skillset <name>`: update named existing skillset, error if not found.

### Delete Version

```bash
asm delete-version --agent <agent_name> --version <version_no>
asm delete-version --skillset <skillset_name> --version <version_no>
```

Uses relative version indexes and accepts version position shifts after deletion.

## Build and Run

```bash
go build ./cmd/asm
./asm --help
```

## Install (Go)

From source repository:

```bash
go install ./cmd/asm
```

Or by module path:

```bash
go install github.com/BaoLe106/asm/cmd/asm@latest
```

## Packaging and Releases

- `.goreleaser.yaml` builds Windows/Linux/macOS binaries.
- Homebrew formula generation is included (configure repository owner/name first).

Release with goreleaser:

```bash
goreleaser release --clean
```

Snapshot build:

```bash
goreleaser release --snapshot --clean
```

## Helper Scripts

- `scripts/build-all.sh`
- `scripts/build-all.ps1`
- `scripts/install.sh`
- `scripts/install.ps1`

## Notes

- This is a starter pack with a placeholder GC implementation for unreferenced object cleanup.
- All state is local to the repository (`.asm`) similar to Git local metadata.
