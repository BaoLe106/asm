#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
go install "$ROOT/cmd/asm"
echo "Installed asm via go install"
