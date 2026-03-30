$ErrorActionPreference = 'Stop'
$Root = (Resolve-Path "$PSScriptRoot/.." ).Path
$Out = Join-Path $Root "dist/local"
New-Item -ItemType Directory -Force -Path $Out | Out-Null

$targets = @(
  @{ GOOS = 'windows'; GOARCH = 'amd64'; Out = 'asm-windows-amd64.exe' },
  @{ GOOS = 'windows'; GOARCH = 'arm64'; Out = 'asm-windows-arm64.exe' },
  @{ GOOS = 'linux'; GOARCH = 'amd64'; Out = 'asm-linux-amd64' },
  @{ GOOS = 'linux'; GOARCH = 'arm64'; Out = 'asm-linux-arm64' },
  @{ GOOS = 'darwin'; GOARCH = 'amd64'; Out = 'asm-darwin-amd64' },
  @{ GOOS = 'darwin'; GOARCH = 'arm64'; Out = 'asm-darwin-arm64' }
)

foreach ($t in $targets) {
  $env:GOOS = $t.GOOS
  $env:GOARCH = $t.GOARCH
  go build -o (Join-Path $Out $t.Out) "$Root/cmd/asm"
}

Remove-Item Env:GOOS -ErrorAction SilentlyContinue
Remove-Item Env:GOARCH -ErrorAction SilentlyContinue
Write-Output "Build output: $Out"
