# cyhi installer - Windows PowerShell 5.1+.
# Usage:
#   .\install.ps1
#   .\install.ps1 -Team "Null Pointers" -Track 3
# If PowerShell blocks the script:
#   powershell -ExecutionPolicy Bypass -File .\install.ps1 -Team "..." -Track 3

[CmdletBinding()]
param(
  [string]$Team,
  [string]$Track,
  [string]$Repo = (Get-Location).Path,
  [string]$SkillsDir = (Join-Path $env:USERPROFILE ".claude\skills"),
  [switch]$NoInit,
  [switch]$NoAgents,
  [switch]$NoPatch
)

$ErrorActionPreference = "Stop"
function Die($m) { Write-Error $m; exit 1 }

# --- locate python ---------------------------------------------------
$py = $null
foreach ($c in @("py -3", "python3", "python")) {
  $parts = $c.Split(" ")
  $exe = Get-Command $parts[0] -ErrorAction SilentlyContinue
  if (-not $exe) { continue }
  try {
    $v = & $parts[0] @($parts[1..($parts.Length-1)]) -c "import sys;print(sys.version_info>=(3,8))" 2>$null
    if ($v -match "True") { $py = $c; break }
  } catch { }
}
if (-not $py) { Die "Python 3.8+ not found. Install from python.org (tick 'Add to PATH'), then rerun." }
$pyExe  = $py.Split(" ")[0]
$pyArgs = @($py.Split(" ") | Select-Object -Skip 1)
function Invoke-Py { param([string[]]$A) & $pyExe @($pyArgs + $A) }

# --- locate the skill source ----------------------------------------
$self = Split-Path -Parent $MyInvocation.MyCommand.Path
$src = $null
foreach ($cand in @(
    (Join-Path $self "cyhi-skills"), $self,
    (Join-Path (Get-Location) "cyhi-skills"), (Get-Location).Path)) {
  if ((Test-Path (Join-Path $cand "SKILL.md")) -and (Test-Path (Join-Path $cand "assets\cyhi"))) {
    $src = $cand; break
  }
}
if (-not $src) { Die "cyhi-skills not found. Run this from the unzipped folder." }

# --- install the skill ----------------------------------------------
$dest = Join-Path $SkillsDir "cyhi-skills"
if (Test-Path $dest) { Remove-Item $dest -Recurse -Force }
New-Item -ItemType Directory -Path (Join-Path $dest "assets") -Force | Out-Null
foreach ($f in @("SKILL.md", "install.md", "tracks.md")) {
  Copy-Item (Join-Path $src $f) (Join-Path $dest $f) -Force
}
foreach ($f in @("cyhi", "AGENTS.md.tmpl")) {
  Copy-Item (Join-Path $src "assets\$f") (Join-Path $dest "assets\$f") -Force
}
Write-Host "installed skill -> $dest"

# --- patch: `cyhi log` blocks forever on an empty, non-tty stdin -----
if (-not $NoPatch) {
  $patch = @'
import sys, pathlib
p = pathlib.Path(sys.argv[1]); s = p.read_text(encoding="utf-8")
old = '"""Append an action record describing what the agent just did."""\n    payload = stdin_json()'
new = '"""Append an action record describing what the agent just did."""\n    payload = stdin_json() if os.environ.get("CYHI_STDIN") else {}'
if old in s:
    p.write_text(s.replace(old, new, 1), encoding="utf-8")
    print("patched: cyhi log no longer blocks on empty stdin")
elif new in s:
    print("patch already applied")
else:
    print("warning: could not apply stdin patch")
'@
  $tmp = Join-Path $env:TEMP "cyhi_patch.py"
  Set-Content -Path $tmp -Value $patch -Encoding UTF8
  Invoke-Py @($tmp, (Join-Path $dest "assets\cyhi"))
  Remove-Item $tmp -Force
}

# --- set up the repo -------------------------------------------------
if ((-not $NoInit) -and ($Team -or $Track)) {
  if (-not $Team)  { Die "-Track given without -Team" }
  if (-not $Track) { Die "-Team given without -Track" }
  if ($Track -notmatch '^[1-4]$') { Die "-Track must be 1, 2, 3 or 4" }
  if (-not (Test-Path $Repo)) { Die "no such directory: $Repo" }

  # auto-init copies into cyhi-logs\bin but does not create it
  New-Item -ItemType Directory -Path (Join-Path $Repo "cyhi-logs\bin") -Force | Out-Null
  Push-Location $Repo
  try {
    Invoke-Py @((Join-Path $dest "assets\cyhi"), "auto-init", "--team", $Team, "--track", $Track)
  } finally { Pop-Location }

  # Claude Code runs hooks through cmd.exe: it cannot expand ${VAR:-.} and will
  # not run an extensionless Python file. Rewrite the three commands.
  $fix = @'
import json, sys, pathlib
path, py = pathlib.Path(sys.argv[1]), sys.argv[2]
cfg = json.loads(path.read_text(encoding="utf-8"))
n = 0
for entries in cfg.get("hooks", {}).values():
    for entry in entries:
        for h in entry.get("hooks", []):
            c = h.get("command", "")
            if "cyhi" in c:
                h["command"] = '%s "%%CLAUDE_PROJECT_DIR%%\\cyhi-logs\\bin\\cyhi" %s' % (py, c.rsplit(" ", 1)[-1])
                n += 1
path.write_text(json.dumps(cfg, indent=2) + "\n", encoding="utf-8")
print("rewrote %d hook command(s) for Windows" % n)
'@
  $tmp2 = Join-Path $env:TEMP "cyhi_fixhooks.py"
  Set-Content -Path $tmp2 -Value $fix -Encoding UTF8
  Invoke-Py @($tmp2, (Join-Path $Repo ".claude\settings.json"), $py)
  Remove-Item $tmp2 -Force

  $agents = Join-Path $Repo "AGENTS.md"
  if ((-not $NoAgents) -and (-not (Test-Path $agents))) {
    (Get-Content (Join-Path $dest "assets\AGENTS.md.tmpl") -Raw).
      Replace("{{TEAM}}", $Team).Replace("{{TRACK}}", $Track) |
      Set-Content -Path $agents -Encoding UTF8 -NoNewline
    $claude = Join-Path $Repo "CLAUDE.md"
    if (-not (Test-Path $claude)) { Set-Content -Path $claude -Value "@AGENTS.md" -Encoding UTF8 }
    Write-Host "wrote AGENTS.md and CLAUDE.md"
  }

  Write-Host ""
  Write-Host "Next:"
  Write-Host "  1. restart your agent so the hooks load"
  Write-Host "  2. $py cyhi-logs\bin\cyhi whoami --set <your roll number>"
  Write-Host "  3. git add -A; git commit -m 'Add cyhi session logging'"
  Write-Host ""
  Write-Host "On Windows, run the tool as: $py cyhi-logs\bin\cyhi <command>"
} else {
  Write-Host ""
  Write-Host "Skill installed. To set up a repo:"
  Write-Host "  cd <repo>; $py `"$dest\assets\cyhi`" auto-init --team `"NAME`" --track N"
}
