$Repo = "thxrsxm/harzmind-code"
$BinaryName = "hzmind.exe"
$InstallDir = "$env:LOCALAPPDATA\HarzMindCode"

$LatestReleaseUrl = "https://api.github.com/repos/$Repo/releases/latest"
$LatestTag = (Invoke-RestMethod -Uri $LatestReleaseUrl).tag_name

$DownloadUrl = "https://github.com/$Repo/releases/download/$LatestTag/hzmind-windows-amd64.exe"

if (!(Test-Path -Path $InstallDir)) {
    New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
}

Write-Host "Donwload version $LatestTag ..."
$OutputPath = Join-Path -Path $InstallDir -ChildPath $BinaryName
Invoke-WebRequest -Uri $DownloadUrl -OutFile $OutputPath

# Add to PATH
$CurrentPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($CurrentPath -notlike "*$InstallDir*") {
    [Environment]::SetEnvironmentVariable("Path", "$CurrentPath;$InstallDir", "User")
    Write-Host "Path has been extended. Please restart the terminal." -ForegroundColor Yellow
}

Write-Host "Installation complete! Open a new terminal and type: hzmind" -ForegroundColor Green
