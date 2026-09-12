# CommandSpeak

**Terminal Velocity — CLI tools for developer workflows.**

CommandSpeak is a natural language CLI tool that saves developers time by translating simple English sentences into complex terminal commands. 

**V2 Improvements**: CommandSpeak now features a flexible keyword-scoring parser, a `--dry-run` flag, Git safety confirmations, colored output, and interactive configuration!

## Features

- **Flexible Natural Language Parsing**: You don't have to memorize exact phrases. Say "push code to vercel", "deploy my project to vercel", or "vercel deploy now"—the parser uses keyword scoring to figure out your intent without any external LLM APIs (100% offline).
- **Dry-Run Mode**: Not sure what a sentence will do? Add the `--dry-run` flag or just include "show me" in your sentence to see the exact shell command that would run without actually executing it.
- **Safety Confirmations**: Before pushing or deploying, the tool checks your `git status`. If you have uncommitted changes, it will warn you and ask for confirmation.
- **Interactive Configuration**: Run `commandspeak change` to interactively customize the exact shell scripts mapped to each intent. 

## Installation

Ensure you have [Go](https://go.dev/) installed.

```bash
# Clone the repository
git clone https://github.com/devchavda2007-web/commandspeak
cd commandspeak

# Download dependencies
go mod tidy

# Build and install the binary globally
go install
```

Make sure your `~/go/bin` directory is in your system's PATH.

## Usage Examples

**1. Deploying a project**
```bash
$ commandspeak "deploy my project to vercel"
Deploying to vercel...
> git add . && git commit -m "update" && git push && vercel --prod
```

**2. Pushing code with a commit message**
```bash
$ commandspeak "push my code with message initial commit"
Pushing code with message: 'initial commit'
> git add . && git commit -m "initial commit" && git push
```

**3. Cleaning old branches**
```bash
$ commandspeak "clean my old branches"
Are you sure you want to clean branches? (y/n)
y
Cleaning branches...
> git branch --merged | grep -v "\*" | grep -v "master" | grep -v "main" | xargs -n 1 git branch -d
```

**4. Check Git Status (New!)**
```bash
$ commandspeak "what is my git status"
> git status
```

**5. Undo last commit (New!)**
```bash
$ commandspeak "undo my last commit"
Are you sure you want to undo your last commit? (y/n)
y
> git reset --soft HEAD~1
```

**6. Dry-Run Mode**
```bash
$ commandspeak "show me how to deploy to vercel"
Dry-run mode. Would execute:
> git add . && git commit -m "update" && git push && vercel --prod
```
*(You can also use the `--dry-run` flag like `commandspeak "deploy" --dry-run`)*

## How to Customize

Run `commandspeak change` to view and edit your command templates interactively. Your preferences are saved locally to `~/.commandspeakrc.json` so your secrets and personal workflows are never tracked in Git!
