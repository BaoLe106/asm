# Build and Run

This document covers local build, run, install, and release commands for asm.

## Prerequisites

1. Go 1.22+
2. A repository folder containing this project

## Build

### Windows (PowerShell)

```powershell
go build -o asm.exe ./cmd/asm
```

### Linux/macOS (bash)

```bash
go build -o asm ./cmd/asm
```

## Run

### Windows

```powershell
.\asm.exe --help
.\asm.exe status
```

### Linux/macOS

```bash
./asm --help
./asm status
```

## Common Command Examples

```bash
asm upsert v1
asm checkout v1
asm list --version
asm list --agent
asm list --skill
asm list --skill --agent-name github
asm delete-version v1
asm status
```

## Run from Source Without Building Binary

```bash
go run ./cmd/asm --help
go run ./cmd/asm list --version
```

## Install

### Local source install

```bash
go install ./cmd/asm
```

### Module install

```bash
go install github.com/BaoLe106/asm/cmd/asm@latest
```
