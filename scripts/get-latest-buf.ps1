# Download the latest buf binary from GitHub Releases for Windows.
# Example: .\get-latest-buf.ps1 -Version v1.34.0 -Destination C:\tools\buf.exe

param(
  [string]$Version = "latest",
  [string]$Destination
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

function Get-BufArchitecture {
  if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64" -or $env:PROCESSOR_ARCHITEW6432 -eq "ARM64") {
    return "arm64"
  }
  if ($env:PROCESSOR_ARCHITECTURE -eq "AMD64" -or $env:PROCESSOR_ARCHITEW6432 -eq "AMD64") {
    return "x86_64"
  }
  throw "Unsupported architecture: $($env:PROCESSOR_ARCHITECTURE)"
}

$arch = Get-BufArchitecture

if (-not $Destination) {
  $Destination = Join-Path $env:USERPROFILE "prj/bin/buf.exe"
}

$directory = Split-Path -Parent $Destination
if (-not (Test-Path -Path $directory)) {
  New-Item -ItemType Directory -Path $directory -Force | Out-Null
}

if ($Version -eq "latest") {
  $downloadUri = "https://github.com/bufbuild/buf/releases/latest/download/buf-Windows-$arch.exe"
  $tag = "latest"
} else {
  $tag = $Version
  if (-not $Version.StartsWith("v")) {
    $tag = "v$Version"
  }
  $downloadUri = "https://github.com/bufbuild/buf/releases/download/$tag/buf-Windows-$arch.exe"
}

Write-Host "Downloading buf ($arch) from $downloadUri ..."
Invoke-WebRequest -Uri $downloadUri -OutFile $Destination -Headers @{ "User-Agent" = "pathist-buf-downloader" } -ErrorAction Stop

$downloaded = Get-Item -Path $Destination
Write-Host "Downloaded version tag: $tag"
Write-Host ("Saved to {0} ({1} bytes)" -f $downloaded.FullName, $downloaded.Length)
