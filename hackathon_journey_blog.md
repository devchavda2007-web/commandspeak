# Building CommandSpeak: A 24-Hour Hackathon Journey

Welcome to the dev blog for **CommandSpeak**, a natural language CLI tool built during a 24-hour hackathon for Track 4: *Terminal Velocity (CLI Tools for Developer Workflows)*.

Our mission was simple: eliminate the friction of terminal context-switching and help developers execute complex commands using plain English.

Here is a look back at our development journey, told through our Git commit history.

---

## Step 1: The Core Engine Revamp
**Commit:** `c596061 feat: V2 improvements (flexible parser, dry-run, colors, tests)`

We started strong by rewriting the core parser. V1 was rigid, so we implemented a **flexible keyword-scoring parser** that runs 100% locally. 
- You can now say "deploy my project" or "vercel deploy now"—the parser figures it out.
- We added a crucial `--dry-run` feature so users can safely check what commands will execute.
- Colored outputs made the terminal UI readable and clean.

## Step 2: Expanding the Vocabulary
**Commit:** `2b88a64 feat: Add INIT, INSTALL, START intents`

A CLI tool is only as good as the workflows it supports. We expanded the parser's vocabulary to understand new critical intents:
- `INIT`: "initialize a new repository" -> `git init`
- `INSTALL`: "install dependencies" -> `npm install`
- `START`: "start the dev server" -> `npm start`

## Step 3: The Big Idea - The Repo Manager
**Commit:** `1dcf76a feat: Add repo manager subcommand and CLONE intent`

We realized that managing remote repositories is a huge pain point. We introduced the `CLONE` intent and the foundations of the `repo` subcommand. This allowed users to simply say:
`"clone repo https://github.com/user/project"`
The tool would dynamically extract the URL and run the git clone command, saving developers from manually copying and pasting URLs.

## Step 4: Frictionless Inputs
**Commit:** `a4ec487 feat: repo command now accepts URL via arg, flag, or prompt`

User experience is everything in CLI design. We wanted to ensure users could launch the repo manager however they liked. We updated the `repo` command to support:
- Direct arguments: `commandspeak repo <url>`
- Flags: `commandspeak repo -u <url>`
- Interactive prompts (asking the user if they didn't provide one).

## Step 5: The Ultimate Interactive Manager & Cross-Platform Polish
**Commit:** `c0f34cb fix: Windows-compatible shell, complete repo manager with 10 options, clean UI`

This was the crown jewel of the hackathon. We turned the `repo` command into a fully interactive workspace.
- **The 10-Option Menu:** Users can now list files, view contents, edit natively, commit, push, pull, check status, check logs, view branches, and switch repos—all via a numbered menu.
- **Windows Compatibility:** We ensured that the underlying `shellRun` function detects the OS. On Windows, edits spawn `notepad` and commands run via `cmd.exe`. On Linux/Mac, it uses `nano` and `bash`.

## Step 6: Hackathon Readiness & Documentation
**Commit:** `a1f4501 docs: update README and add design note for Track 4 submission`
**Commit:** `c360ed7 Add cyhi session logging and update blog`

As the clock ticked down, we initialized the CYHI hackathon logging for "Team eteranel byte". We finalized the mandatory `design_note.md` required by the track guidelines and updated the `README.md` to showcase the new interactive repo manager.

---
*Built with ❤️ by Eternal Byte for the CYHI Hackathon.*
