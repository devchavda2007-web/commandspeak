# Design Note: CommandSpeak

## The Problem
Developers waste significant time performing repetitive terminal tasks: setting up projects, managing git workflows, switching environments, and looking up long CLI flags. While aliases exist, they are hard to share across teams and quickly become unmaintainable. Switching contexts between the browser, IDE, and terminal for simple tasks like editing a file in a remote repo breaks focus.

## Key Design Choices
1. **Natural Language via Keyword Scoring (100% Offline)**
   Instead of using an LLM API which introduces latency, API keys, and internet dependence, CommandSpeak uses a lightweight, offline keyword-scoring algorithm. It's fast, robust, and translates intents accurately (e.g., "show me my git status" vs "deploy to vercel").

2. **The `repo` Interactive Manager**
   To solve the context-switching problem, we built an interactive repo manager right into the CLI. Users can clone, browse, view, edit (spawning OS-native editors like Notepad/nano), commit, and push without ever touching an external GUI.

3. **User Configurable Intents (`~/.commandspeakrc.json`)**
   We recognized that "deploy" means `vercel --prod` to one developer, and `docker-compose up -d` to another. The `commandspeak change` interactive prompt lets users overwrite the default templates. These are stored securely in the user's home directory so they never accidentally leak into version control.

4. **Safety & Transparency (`--dry-run`)**
   Translating English to bash can be risky. We implemented a robust `--dry-run` (or "show me") feature. Additionally, destructive operations (cleaning branches, undoing commits, deploying uncommitted code) require explicit `(y/n)` confirmations and automatically check `git status`.

## Limitations & Future Work
- **Static Keyword Mapping**: While fast, the keyword-scoring approach struggles with complex compound commands (e.g., "deploy to vercel but first run my tests"). Future versions could implement a small, locally-running NLP model to parse abstract syntax trees from English phrases.
- **Hardcoded Placeholder Extraction**: Parameter extraction (like URLs or platform names) currently relies on specific Regex patterns. A more generic templating engine could allow users to define their own extracted variables.
