# wslbm updater
$ErrorActionPreference = "Stop"

$installDir = "$env:LOCALAPPDATA\wslbm"
$distrosDir = "$installDir\distros"
$zipUrl     = "https://github.com/Catergems/wslbm/releases/download/Tags/wslbm.zip"
$zipPath    = "$env:TEMP\wslbm-update.zip"
$extractDir = "$env:TEMP\wslbm-update-extract"

Write-Host "Updating wslbm..."

# Kill any running wslbm processes
$procs = Get-Process -Name "wslbm" -ErrorAction SilentlyContinue
if ($procs) {
    Write-Host "Terminating running wslbm instances..."
    $procs | Stop-Process -Force
    Start-Sleep -Milliseconds 500
}

Write-Host "Downloading latest release..."
curl -L $zipUrl -o $zipPath

Write-Host "Extracting..."
if (Test-Path $extractDir) { Remove-Item $extractDir -Recurse -Force }
Expand-Archive -Path $zipPath -DestinationPath $extractDir -Force

$exe = Get-ChildItem -Path $extractDir -Filter "wslbm.exe" -Recurse | Select-Object -First 1
if (-not $exe) { Write-Error "wslbm.exe not found in zip."; exit 1 }

$oldExe = "$installDir\wslbm.old.exe"
if (Test-Path "$installDir\wslbm.exe") { Move-Item "$installDir\wslbm.exe" $oldExe -Force }

try {
    Copy-Item $exe.FullName -Destination "$installDir\wslbm.exe" -Force
    if (Test-Path $oldExe) { Remove-Item $oldExe -Force }
} catch {
    if (Test-Path $oldExe) { Move-Item $oldExe "$installDir\wslbm.exe" -Force }
    Write-Error "Update failed: $_"
    exit 1
}

$distrosSrc = Get-ChildItem -Path $extractDir -Filter "distros" -Recurse -Directory | Select-Object -First 1
if ($distrosSrc) {
    Write-Host "Updating distro definitions..."
    Copy-Item "$($distrosSrc.FullName)\*" -Destination $distrosDir -Force
}

Remove-Item $zipPath -Force
Remove-Item $extractDir -Recurse -Force

Write-Host ""
Write-Host "wslbm updated successfully!"
