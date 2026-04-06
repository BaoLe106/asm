# asm

AI Agent Skills Manager CLI written in Go.

asm manages local AI agent setups in a repository using a Git-like hidden store (.asm), with version-name based snapshots.

## Key Concepts

- Working tree: .{agent}/skills/{skillName} plus any additional files under .{agent}
- Local metadata/object store: .asm/
- Version model:

1. Each version is identified by a string name (example: v1, baseline, release-2026-04)
2. Each version stores all current agents and their snapshot IDs
3. One version is marked as current in .asm/state/current.json

## Local Storage Layout

```text
.asm/
  state/current.json
  objects/<hash-prefix>/<hash>
  snapshots/<snapshot-id>.json
  refs/versions/<version_name>/version.json
  locks/write.lock
```

## Commands

### checkout

Checkout a version into the working tree.

```bash
asm checkout <version>
```

### upsert

Create or update a named version using the current local agent folders.

```bash
asm upsert <version>
```

### delete-version

Delete a named version.

```bash
asm delete-version <version>
```

### list

List versions or list agents/skills from the current version.

```bash
asm list --version
asm list --agent
asm list --skill
asm list --skill --agent-name <agent>
```

Behavior:

1. asm list --version: lists all version names
2. asm list --agent: lists all agents in the current version
3. asm list --skill: lists all skills for all agents in the current version as agent:skill
4. asm list --skill --agent-name <agent>: lists only skills for one agent from the current version

### status

Show current version.

```bash
asm status
```

## Build and Run

Build/run instructions are in [BUILD_AND_RUN.md](BUILD_AND_RUN.md).

## Notes

1. All state is local to the repository (.asm), similar to Git local metadata.
2. Hidden folder behavior depends on OS and file attributes.
