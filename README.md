<div align="center">
  <h1>🚀 CommandSpeak</h1>
  <p><b>Terminal Velocity — Natural Language CLI Tools for Developer Workflows</b></p>
  <p><i>A Track 4 Submission for CYHI Hackathon (Team: eteranel byte)</i></p>
</div>

---

**CommandSpeak** is a powerful CLI tool that saves developers time by translating simple English sentences into complex terminal commands. 

No more googling obscure git flags or memorizing exact bash scripts. Just speak plain English, and CommandSpeak executes the right developer commands. Best of all, it works **100% offline** with no external LLM API keys required.

## ✨ Features

- 🧠 **Natural Language Parsing:** Say `"push code to vercel"`, `"deploy my project"`, or `"clean my old branches"`. The built-in keyword scoring engine instantly figures out your intent.
- 📦 **Interactive Repo Manager:** Run `commandspeak repo` to launch a fully interactive TUI (Terminal User Interface). Browse files, edit in Nano/Notepad, commit, push, pull, and branch—without leaving the menu.
- 🖥️ **Private Web Dashboard:** Run `commandspeak ui` to pop open a sleek web UI showing your command history and settings. Data is stored safely in a **local SQLite database** and never leaves your machine.
- 🛡️ **Git Safety Confirmations:** Before pushing or deploying, it checks your `git status`. If you have uncommitted changes, it warns you before making a mess.
- 🔍 **Dry-Run Mode:** Add `--dry-run` or say `"show me how to deploy"` to see exactly what bash command would run without actually executing it.
- ⚙️ **Fully Customizable:** Run `commandspeak change` to edit your command templates interactively.

---

## 🚀 Installation

Ensure you have [Go](https://go.dev/) installed on your machine.

```bash
# 1. Clone the repository
git clone https://github.com/devchavda2007-web/commandspeak.git
cd commandspeak

# 2. Download Go dependencies
go mod tidy

# 3. Build and install the CLI globally
go install
```

> **Note:** Ensure your `~/go/bin` directory (or `%USERPROFILE%\go\bin` on Windows) is in your system's `PATH`.

*(Optional)* **Install PHP** for the private SQLite Dashboard:
- Windows: `choco install php`
- Mac: `brew install php`

---

## 📖 Usage Guide

### 1. The Natural Language CLI
Just type `commandspeak` followed by what you want to do:

```bash
$ commandspeak "deploy my project to vercel"
✔ Done! > git add . && git commit -m "update" && git push && vercel --prod

$ commandspeak "show me my git status"
✔ Done! > git status

$ commandspeak "install packages and start dev server"
✔ Done! > npm install && npm start
```

### 2. The Interactive Repo Manager
Want to manage a GitHub repo interactively?

```bash
$ commandspeak repo https://github.com/devchavda2007-web/commandspeak
```
*This opens a 10-option interactive menu to view files, edit natively, commit, push, and more.*

### 3. The Web Dashboard (`commandspeak ui`)
Want to see your command history or configure your settings in a beautiful GUI?

```bash
$ commandspeak ui
```
*This automatically starts a local server and opens the dashboard in your browser.*

### 4. Customizing Commands
Want to change what command runs when you say "deploy"?

```bash
$ commandspeak change
```
*This opens an interactive wizard to edit your `~/.commandspeakrc.json` file.*

---

## 🔒 Privacy & Architecture

CommandSpeak is designed for **maximum privacy**:
1. **Offline NLP**: The natural language parser uses local keyword-scoring heuristics in Go. No code or prompts are sent to OpenAI/Anthropic.
2. **Local SQLite Backend**: The `commandspeak ui` dashboard uses a lightweight PHP server (`api.php`) to save your data to a `private_local_data.sqlite` database stored locally on your machine.
3. **No Cloud Syncing**: Your history, configs, and repos never leave your computer. The database is explicitly git-ignored.

## 🏆 CYHI Hackathon Compliance
- **Track 4 (Terminal Velocity)**: Designed specifically to accelerate CLI workflows.
- All development was logged using the `cyhi log` tool (`cyhi-logs/turns/devchavda2007.jsonl`).
- Implemented core deliverables: The Go CLI, interactive GitHub manager, and private SQLite data retention.
