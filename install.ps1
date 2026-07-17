# wslbm installer
$ErrorActionPreference = "Stop"

$installDir = "$env:LOCALAPPDATA\wslbm"
$distrosDir = "$installDir\distros"
$wslosDir   = "$installDir\wslos"
$zipUrl     = "https://github.com/Catergems/wslbm/releases/download/Tags/wslbm.zip"
$zipPath    = "$env:TEMP\wslbm-install.zip"
$extractDir = "$env:TEMP\wslbm-extract"

Write-Host "Installing wslbm to $installDir..."

New-Item -ItemType Directory -Force -Path $installDir | Out-Null
New-Item -ItemType Directory -Force -Path $distrosDir | Out-Null
New-Item -ItemType Directory -Force -Path $wslosDir   | Out-Null

Write-Host "Downloading wslbm.zip..."
curl -L $zipUrl -o $zipPath

Write-Host "Extracting..."
if (Test-Path $extractDir) { Remove-Item $extractDir -Recurse -Force }
Expand-Archive -Path $zipPath -DestinationPath $extractDir -Force

$exe = Get-ChildItem -Path $extractDir -Filter "wslbm.exe" -Recurse | Select-Object -First 1
if (-not $exe) { Write-Error "wslbm.exe not found in zip."; exit 1 }
Copy-Item $exe.FullName -Destination "$installDir\wslbm.exe" -Force

$distrosSrc = Get-ChildItem -Path $extractDir -Filter "distros" -Recurse -Directory | Select-Object -First 1
if ($distrosSrc) { Copy-Item "$($distrosSrc.FullName)\*" -Destination $distrosDir -Force }

Remove-Item $zipPath -Force
Remove-Item $extractDir -Recurse -Force

$currentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
if ($currentPath -notlike "*$installDir*") {
    Write-Host "Adding $installDir to PATH..."
    [Environment]::SetEnvironmentVariable("PATH", "$currentPath;$installDir", "User")
}

Write-Host ""
Write-Host "wslbm installed! Restart your terminal and run: wslbm help"
