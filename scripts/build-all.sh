#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
OUT="$ROOT/dist/local"
mkdir -p "$OUT"

GOOS=windows GOARCH=amd64 go build -o "$OUT/asm-windows-amd64.exe" "$ROOT/cmd/asm"
GOOS=windows GOARCH=arm64 go build -o "$OUT/asm-windows-arm64.exe" "$ROOT/cmd/asm"
GOOS=linux GOARCH=amd64 go build -o "$OUT/asm-linux-amd64" "$ROOT/cmd/asm"
GOOS=linux GOARCH=arm64 go build -o "$OUT/asm-linux-arm64" "$ROOT/cmd/asm"
GOOS=darwin GOARCH=amd64 go build -o "$OUT/asm-darwin-amd64" "$ROOT/cmd/asm"
GOOS=darwin GOARCH=arm64 go build -o "$OUT/asm-darwin-arm64" "$ROOT/cmd/asm"

echo "Build output: $OUT"
