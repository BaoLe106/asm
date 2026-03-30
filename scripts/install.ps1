$ErrorActionPreference = 'Stop'
$Root = (Resolve-Path "$PSScriptRoot/.." ).Path
go install "$Root/cmd/asm"
Write-Output "Installed asm via go install"
