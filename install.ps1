# wslbm installer
$ErrorActionPreference = "Stop"

$installDir = "$env:LOCALAPPDATA\wslbm"
$distrosDir = "$installDir\distros"
$wslosDir   = "$installDir\wslos"
$versionUrl = "https://raw.githubusercontent.com/Catergems/wslbm/main/version.txt"
$latestVersion = (Invoke-RestMethod -Uri $versionUrl).Trim()
$zipUrl     = "https://github.com/Catergems/wslbm/releases/download/release-wslbm/wslbm-$latestVersion.zip"
$zipPath    = "$env:TEMP\wslbm-install.zip"
$extractDir = "$env:TEMP\wslbm-extract"

Write-Host "Installing wslbm to $installDir..."

# Create directories
New-Item -ItemType Directory -Force -Path $installDir | Out-Null
New-Item -ItemType Directory -Force -Path $distrosDir | Out-Null
New-Item -ItemType Directory -Force -Path $wslosDir   | Out-Null

# Download zip
Write-Host "Downloading wslbm-$latestVersion.zip..."
Invoke-WebRequest -Uri $zipUrl -OutFile $zipPath

# Extract
Write-Host "Extracting..."
if (Test-Path $extractDir) { Remove-Item $extractDir -Recurse -Force }
Expand-Archive -Path $zipPath -DestinationPath $extractDir -Force

# Copy wslbm.exe
$exe = Get-ChildItem -Path $extractDir -Filter "wslbm.exe" -Recurse | Select-Object -First 1
if (-not $exe) {
    Write-Error "wslbm.exe not found in zip."
    exit 1
}
Copy-Item $exe.FullName -Destination "$installDir\wslbm.exe" -Force

# Copy distros folder
$distrosSrc = Get-ChildItem -Path $extractDir -Filter "distros" -Recurse -Directory | Select-Object -First 1
if ($distrosSrc) {
    Copy-Item "$($distrosSrc.FullName)\*" -Destination $distrosDir -Force
}

# Cleanup
Remove-Item $zipPath -Force
Remove-Item $extractDir -Recurse -Force

# Add to PATH if not already present
$currentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
if ($currentPath -notlike "*$installDir*") {
    Write-Host "Adding $installDir to PATH..."
    [Environment]::SetEnvironmentVariable("PATH", "$currentPath;$installDir", "User")
}

Write-Host ""
Write-Host "wslbm installed successfully!"
Write-Host "Restart your terminal and run: wslbm help"
