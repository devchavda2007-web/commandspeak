# CommandSpeak

**Terminal Velocity — CLI tools for developer workflows.**

CommandSpeak is a natural language CLI tool that saves developers time by translating simple English sentences into complex terminal commands. 

**V2 Improvements**: CommandSpeak now features a flexible keyword-scoring parser, a `--dry-run` flag, Git safety confirmations, colored output, interactive configuration, and a built-in interactive GitHub Repository Manager!

## Features

- **Interactive GitHub Repo Manager**: Manage any GitHub repository right from the CLI. Clone, view files, edit natively, commit, and pull via an interactive menu. Run `commandspeak repo`.
- **Flexible Natural Language Parsing**: You don't have to memorize exact phrases. Say "push code to vercel", "deploy my project", or "clone repo https://github..."—the parser uses keyword scoring to figure out your intent without any external LLM APIs (100% offline).
- **Dry-Run Mode**: Not sure what a sentence will do? Add the `--dry-run` flag or just include "show me" in your sentence to see the exact shell command that would run without actually executing it.
- **Safety Confirmations**: Before pushing, deploying, or undoing commits, the tool checks your `git status`. If you have uncommitted changes, it will warn you and ask for confirmation.
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

**1. Interactive GitHub Repo Manager (New!)**
```bash
# Interactive mode
$ commandspeak repo

# Or pass a URL directly
$ commandspeak repo https://github.com/user/project
```
This drops you into an interactive menu with 10 options allowing you to browse files, edit files natively (Notepad/nano), push commits, and more.

**2. Natural Language Cloning**
```bash
$ commandspeak "clone repo https://github.com/user/project"
> git clone https://github.com/user/project
```

**3. Deploying a project**
```bash
$ commandspeak "deploy my project to vercel"
> git add . && git commit -m "update" && git push && vercel --prod
```

**4. Check Git Status**
```bash
$ commandspeak "what is my git status"
> git status
```

**5. Start Server / Install Dependencies**
```bash
$ commandspeak "install packages and start dev server"
> npm install && npm start
```

**6. Dry-Run Mode**
```bash
$ commandspeak "show me how to deploy to vercel"
Dry-run — would execute:
  > git add . && git commit -m "update" && git push && vercel --prod
```

## How to Customize

Run `commandspeak change` to view and edit your command templates interactively. Your preferences are saved locally to `~/.commandspeakrc.json` so your personal workflows are never tracked in Git!
