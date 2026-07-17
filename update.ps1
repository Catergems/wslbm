# wslbm updater
param(
    [string]$LatestVersion
)

$ErrorActionPreference = "Stop"

# Wait 1 second to give the calling wslbm process time to exit and release file locks
Start-Sleep -Seconds 1

$installDir = "$env:LOCALAPPDATA\wslbm"
$distrosDir = "$installDir\distros"

if (-not $LatestVersion) {
    Write-Host "Checking latest version..."
    $versionUrl = "https://raw.githubusercontent.com/Catergems/wslbm/main/version.txt"
    $LatestVersion = (Invoke-RestMethod -Uri $versionUrl).Trim()
}

$zipUrl     = "https://github.com/Catergems/wslbm/releases/download/release-wslbm/wslbm-$LatestVersion.zip"
$zipPath    = "$env:TEMP\wslbm-update.zip"
$extractDir = "$env:TEMP\wslbm-update-extract"

Write-Host "Updating wslbm to version $LatestVersion..."

# Download zip
Write-Host "Downloading latest release from $zipUrl..."
Invoke-WebRequest -Uri $zipUrl -OutFile $zipPath

# Extract
Write-Host "Extracting..."
if (Test-Path $extractDir) { Remove-Item $extractDir -Recurse -Force }
Expand-Archive -Path $zipPath -DestinationPath $extractDir -Force

# Replace wslbm.exe
$exe = Get-ChildItem -Path $extractDir -Filter "wslbm.exe" -Recurse | Select-Object -First 1
if (-not $exe) {
    Write-Error "wslbm.exe not found in zip."
    exit 1
}

# Rename old exe before replacing
$oldExe = "$installDir\wslbm.old.exe"
if (Test-Path "$installDir\wslbm.exe") {
    Move-Item "$installDir\wslbm.exe" $oldExe -Force
}

try {
    Copy-Item $exe.FullName -Destination "$installDir\wslbm.exe" -Force
    
    # Retry removing the old exe in case the process is still releasing
    $attempts = 0
    while ($attempts -lt 5) {
        try {
            if (Test-Path $oldExe) {
                Remove-Item $oldExe -Force -ErrorAction Stop
            }
            break
        } catch {
            Start-Sleep -Seconds 1
            $attempts++
        }
    }
} catch {
    # Restore old exe on failure
    if (Test-Path $oldExe) { Move-Item $oldExe "$installDir\wslbm.exe" -Force }
    Write-Error "Update failed: $_"
    exit 1
}

# Update distros folder (preserve wslos/)
$distrosSrc = Get-ChildItem -Path $extractDir -Filter "distros" -Recurse -Directory | Select-Object -First 1
if ($distrosSrc) {
    Write-Host "Updating distro definitions..."
    Copy-Item "$($distrosSrc.FullName)\*" -Destination $distrosDir -Force
}

# Cleanup
Remove-Item $zipPath -Force
Remove-Item $extractDir -Recurse -Force

Write-Host ""
Write-Host "wslbm updated successfully to $LatestVersion!"
Start-Sleep -Seconds 3
