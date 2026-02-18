$Repo = "thxrsxm/harzmind-code"
$BinaryName = "hzmind.exe"
$DownloadName = "hzmind-windows-amd64.exe"
$InstallDir = "$env:LOCALAPPDATA\HarzMindCode"

$LatestReleaseUrl = "https://api.github.com/repos/$Repo/releases/latest"

try {
    $LatestTag = (Invoke-RestMethod -Uri $LatestReleaseUrl).tag_name
} catch {
    Write-Host "Failed to fetch latest release: $_" -ForegroundColor Red
    exit 1
}

$DownloadUrl = "https://github.com/$Repo/releases/download/$LatestTag/$DownloadName"

if (!(Test-Path -Path $InstallDir)) {
    New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
}

Write-Host "Downloading version $LatestTag ..." -ForegroundColor Cyan

# Download with original name first
$TempPath = Join-Path -Path $InstallDir -ChildPath $DownloadName

try {
    Invoke-WebRequest -Uri $DownloadUrl -OutFile $TempPath -UseBasicParsing
} catch {
    Write-Host "Download failed: $_" -ForegroundColor Red
    exit 1
}

# Rename to hzmind.exe so user can just type "hzmind"
$FinalPath = Join-Path -Path $InstallDir -ChildPath $BinaryName
if (Test-Path $FinalPath) {
    Remove-Item $FinalPath -Force
}
Rename-Item -Path $TempPath -NewName $BinaryName

Write-Host "Download complete." -ForegroundColor Cyan

# Add to PATH (persistent)
$CurrentPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($CurrentPath -notlike "*$InstallDir*") {
    [Environment]::SetEnvironmentVariable("Path", "$CurrentPath;$InstallDir", "User")
    # Also update current session
    $env:Path += ";$InstallDir"
    Write-Host "PATH updated." -ForegroundColor Yellow
}

Write-Host "Installation complete! Open a new terminal and type: hzmind" -ForegroundColor Green
