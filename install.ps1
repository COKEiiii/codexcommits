$ErrorActionPreference = "Stop"

$repo = "COKEiiii/codexcommits"
$arch = if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") { "arm64" } else { "x86_64" }
$assetName = "codexcommits_Windows_$arch.zip"
$release = Invoke-RestMethod "https://api.github.com/repos/$repo/releases/latest"
$asset = $release.assets | Where-Object { $_.name -eq $assetName } | Select-Object -First 1

if (-not $asset) {
    throw "Release asset $assetName was not found."
}

$installDir = Join-Path $env:LOCALAPPDATA "Programs\codexcommits"
$tempDir = Join-Path ([System.IO.Path]::GetTempPath()) ("codexcommits-" + [guid]::NewGuid())
$archive = Join-Path $tempDir $assetName
$checksums = Join-Path $tempDir "checksums.txt"

try {
    New-Item -ItemType Directory -Force $tempDir | Out-Null
    New-Item -ItemType Directory -Force $installDir | Out-Null
    Invoke-WebRequest $asset.browser_download_url -OutFile $archive
    $checksumAsset = $release.assets | Where-Object { $_.name -eq "checksums.txt" } | Select-Object -First 1
    if (-not $checksumAsset) { throw "Release checksums were not found." }
    Invoke-WebRequest $checksumAsset.browser_download_url -OutFile $checksums
    $expected = (Get-Content $checksums | Where-Object { $_ -match "\s+$([regex]::Escape($assetName))$" } | Select-Object -First 1) -split "\s+" | Select-Object -First 1
    $actual = (Get-FileHash $archive -Algorithm SHA256).Hash
    if (-not $expected -or $actual -ne $expected) { throw "Release checksum verification failed." }
    Expand-Archive $archive -DestinationPath $tempDir -Force
    Copy-Item (Join-Path $tempDir "codexcommits.exe") $installDir -Force
} finally {
    if (Test-Path $tempDir) { Remove-Item $tempDir -Recurse -Force }
}

$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
$parts = @($userPath -split ";" | Where-Object { $_ })
if ($parts -notcontains $installDir) {
    $newPath = (($parts + $installDir) -join ";")
    [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
}
if (($env:Path -split ";") -notcontains $installDir) {
    $env:Path = "$env:Path;$installDir"
}

Write-Host "Installed codexcommits $($release.tag_name) to $installDir"
Write-Host "Run: codexcommits --version"
