# Setting up cyhi in a project

**NEW: Auto-initialization on first skill call.** Manual setup below is for reference only; the skill handles it automatically.

## Automatic Setup (Recommended)

First time working on this project:

```bash
cyhi-logs/bin/cyhi auto-init --team "TEAM NAME" --track N
```

This single command:
1. Creates `cyhi-logs/`, `bin/`, `turns/`, `.local/` directories
2. Installs hooks to `.claude/settings.json` (prompts auto-captured)
3. Sets `.gitattributes` and `.gitignore`
4. Initializes `state.json` and `HANDOFF.md`
5. Prints handoff for the session

No additional steps. Hooks installed automatically; every prompt is logged.

## Manual Setup (Reference)

If you prefer manual control:

```bash
mkdir -p cyhi-logs/bin
cp ~/.claude/skills/cyhi-skills/assets/cyhi cyhi-logs/bin/cyhi
chmod +x cyhi-logs/bin/cyhi
cyhi-logs/bin/cyhi init --team "TEAM NAME" --track N
cyhi-logs/bin/cyhi install-hooks
```

Then commit `.claude/settings.json`.

This layer is what makes the record survive 24 hours. Without it the model has to remember
to transcribe every prompt by hand, and it will drift.

## 3. Install the portable contract

```bash
cp ~/.claude/skills/cyhi-skills/assets/AGENTS.md.tmpl AGENTS.md
printf '@AGENTS.md\n' > CLAUDE.md
```

Edit `AGENTS.md` to fill in the team name and track. This is what Codex, Cursor, Gemini CLI,
Amp, opencode and Copilot read - none of them load Claude Code skills, but all of them read
a markdown file at the repo root.

## 4. Commit

```bash
git add -A && git commit -m "Add cyhi session logging"
```

---

## Identity

Resolved automatically, in order: `$CYHI_MEMBER` → `cyhi-logs/.local/member` (written once
per machine) → `git config user.email`. Normalised to lowercase.

Pin a roll number on this machine:
```bash
cyhi-logs/bin/cyhi whoami --set 23BCE1234
```
Two people sharing one laptop:
```bash
CYHI_MEMBER=23BCE9999 claude
```

## Layout

| Path | Committed | Purpose |
|---|---|---|
| `cyhi-logs/turns/<member>.jsonl` | yes | append-only record, one file per member |
| `cyhi-logs/state.json` | yes | team, track, roster |
| `cyhi-logs/session.md` | yes | generated report for organisers |
| `cyhi-logs/bin/cyhi` | yes | the tool |
| `cyhi-logs/.local/` | **no** | member id, session id, turn counters |
| `HANDOFF.md` | yes | snapshot for the next session |

One file per member means two people on two laptops never touch the same file. The
`.gitattributes` rule `cyhi-logs/**/*.jsonl merge=union` handles the remaining case (same
person, two branches) by keeping both sets of lines instead of raising a conflict.

## JSONL schema

Two record kinds, joined on `(member, session, turn)`:

```json
{"kind":"prompt","ts":"2026-09-12T11:59:36+05:30","session":"sess01","turn":7,"member":"23bce1234","track":"3","text":"<verbatim prompt>"}
{"kind":"action","ts":"2026-09-12T11:59:41+05:30","session":"sess01","turn":7,"member":"23bce1234","track":"3","type":"code","summary":"...","files":["train.py"]}
```

A prompt line with no matching action line means the agent answered without logging - useful
signal, not an error. Never rewrite these files; append only.

## Commands

```
cyhi init      cyhi log       cyhi handoff   cyhi render
cyhi status    cyhi whoami    cyhi install-hooks
```
Run any of them with `--help`.

## Organiser queries

```bash
cyhi-logs/bin/cyhi render                                    # human-readable report
cat cyhi-logs/turns/*.jsonl | jq -r 'select(.kind=="action").type' | sort | uniq -c
cat cyhi-logs/turns/*.jsonl | jq -r 'select(.kind=="prompt").member' | sort | uniq -c
```
