#!/usr/bin/env bash
# cyhi installer - macOS, Linux, WSL, Git Bash.
# Usage:
#   ./install.sh                                 # install the skill only
#   ./install.sh --team "Null Pointers" --track 3   # install + set up the current repo
set -euo pipefail

TEAM=""; TRACK=""; REPO="$PWD"; SKILLS_DIR="${HOME}/.claude/skills"; DO_INIT=1; DO_AGENTS=1; DO_PATCH=1

die() { printf 'error: %s\n' "$1" >&2; exit 1; }
say() { printf '%s\n' "$1"; }

while [ $# -gt 0 ]; do
  case "$1" in
    --team)       TEAM="${2:-}"; shift 2 ;;
    --track)      TRACK="${2:-}"; shift 2 ;;
    --repo)       REPO="${2:-}"; shift 2 ;;
    --skills-dir) SKILLS_DIR="${2:-}"; shift 2 ;;
    --no-init)    DO_INIT=0; shift ;;
    --no-agents)  DO_AGENTS=0; shift ;;
    --no-patch)   DO_PATCH=0; shift ;;
    -h|--help)
      sed -n '2,5p' "$0" | sed 's/^# \{0,1\}//'; exit 0 ;;
    *) die "unknown option: $1" ;;
  esac
done

# --- locate python ---------------------------------------------------
PY=""
for c in python3 python; do
  if command -v "$c" >/dev/null 2>&1; then
    if "$c" -c 'import sys; sys.exit(0 if sys.version_info >= (3,8) else 1)' 2>/dev/null; then
      PY="$c"; break
    fi
  fi
done
[ -n "$PY" ] || die "python 3.8+ not found. Install it, then rerun."

# --- locate the skill source ----------------------------------------
SELF_DIR="$(cd "$(dirname "$0")" && pwd)"
SRC=""
for cand in "$SELF_DIR/cyhi-skills" "$SELF_DIR" "$PWD/cyhi-skills" "$PWD"; do
  if [ -f "$cand/SKILL.md" ] && [ -f "$cand/assets/cyhi" ]; then SRC="$cand"; break; fi
done
[ -n "$SRC" ] || die "cyhi-skills not found. Run this from the unzipped folder."

# --- install the skill ----------------------------------------------
DEST="$SKILLS_DIR/cyhi-skills"
mkdir -p "$SKILLS_DIR"
rm -rf "$DEST"
mkdir -p "$DEST/assets"
cp "$SRC/SKILL.md" "$SRC/install.md" "$SRC/tracks.md" "$DEST/"
cp "$SRC/assets/cyhi" "$SRC/assets/AGENTS.md.tmpl" "$DEST/assets/"
chmod +x "$DEST/assets/cyhi"
say "installed skill -> $DEST"

# --- patch: `cyhi log` blocks forever when stdin is a pipe with no data,
#     which is exactly how an agent's shell tool invokes it. -------------
if [ "$DO_PATCH" = "1" ]; then
  "$PY" - "$DEST/assets/cyhi" <<'PATCH'
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
    print("warning: could not apply stdin patch; run `cyhi log ... < /dev/null`", file=sys.stderr)
PATCH
fi

# --- set up the repo -------------------------------------------------
if [ "$DO_INIT" = "1" ] && [ -n "$TEAM$TRACK" ]; then
  [ -n "$TEAM" ]  || die "--track given without --team"
  [ -n "$TRACK" ] || die "--team given without --track"
  case "$TRACK" in 1|2|3|4) ;; *) die "--track must be 1, 2, 3 or 4" ;; esac
  [ -d "$REPO" ] || die "no such directory: $REPO"

  # auto-init copies into cyhi-logs/bin/ but does not create it
  mkdir -p "$REPO/cyhi-logs/bin"
  ( cd "$REPO" && "$PY" "$DEST/assets/cyhi" auto-init --team "$TEAM" --track "$TRACK" )

  # Git Bash / MSYS: Claude Code runs hooks through cmd.exe, which cannot expand
  # ${CLAUDE_PROJECT_DIR:-.} and will not resolve an extensionless script.
  case "$(uname -s 2>/dev/null || echo unknown)" in
    MINGW*|MSYS*|CYGWIN*)
      "$PY" - "$REPO/.claude/settings.json" "$PY" <<'FIXHOOKS'
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
FIXHOOKS
      ;;
  esac

  if [ "$DO_AGENTS" = "1" ] && [ ! -f "$REPO/AGENTS.md" ]; then
    sed -e "s/{{TEAM}}/$TEAM/g" -e "s/{{TRACK}}/$TRACK/g" \
      "$DEST/assets/AGENTS.md.tmpl" > "$REPO/AGENTS.md"
    [ -f "$REPO/CLAUDE.md" ] || printf '@AGENTS.md\n' > "$REPO/CLAUDE.md"
    say "wrote AGENTS.md and CLAUDE.md"
  fi

  say ""
  say "Next:"
  say "  1. restart your agent so the hooks load"
  say "  2. cyhi-logs/bin/cyhi whoami --set <your roll number>"
  say "  3. git add -A && git commit -m 'Add cyhi session logging'"
else
  say ""
  say "Skill installed. To set up a repo:"
  say "  cd <repo> && $DEST/assets/cyhi auto-init --team \"NAME\" --track N"
fi
