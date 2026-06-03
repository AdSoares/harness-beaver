# Build do HarnessBeaver (bvr).
#   .\build.ps1         -> compila bvr.exe para a plataforma atual
#   .\build.ps1 -All    -> cross-compila para win/linux/mac em ./dist
param([switch]$All)

$ErrorActionPreference = "Stop"

# Versão: tag/descrição do git, ou data se não houver git.
$version = (git describe --tags --always --dirty 2>$null)
if (-not $version) { $version = "v0.0.0+$(Get-Date -Format yyyyMMdd)" }
$ldflags = "-s -w -X harnessbeaver/cmd.Version=$version"

go build -ldflags $ldflags -o bvr.exe .
Write-Host "Compilado: bvr.exe" -ForegroundColor Green

if ($All) {
    New-Item -ItemType Directory -Force dist | Out-Null
    $targets = @(
        @{ os = "windows"; arch = "amd64"; out = "bvr-windows-amd64.exe" },
        @{ os = "linux";   arch = "amd64"; out = "bvr-linux-amd64" },
        @{ os = "darwin";  arch = "arm64"; out = "bvr-darwin-arm64" }
    )
    foreach ($t in $targets) {
        $env:GOOS = $t.os
        $env:GOARCH = $t.arch
        go build -ldflags $ldflags -o "dist/$($t.out)" .
        Write-Host "dist/$($t.out)" -ForegroundColor Green
    }
    Remove-Item Env:GOOS, Env:GOARCH -ErrorAction SilentlyContinue
}
