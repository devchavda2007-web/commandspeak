# CommandSpeak

**Terminal Velocity — CLI tools for developer workflows.**

CommandSpeak is a natural language CLI tool that saves developers time by translating simple English sentences into complex terminal commands for common tasks like deploying, pushing code, cleaning branches, and running tests.

## Features

- **Natural Language Parsing**: Just type what you want to do. No LLM dependency—runs entirely offline!
- **Interactive Mode**: Run `commandspeak` without arguments to enter an interactive shell.
- **Single Command Mode**: Run `commandspeak "your sentence here"` to execute a command instantly.
- **Configurable**: Automatically loads your preferences from `~/.commandspeakrc.json`.

## Installation

Ensure you have [Go](https://go.dev/) installed.

```bash
# Clone the repository (or copy the folder)
cd commandspeak

# Build and install the binary globally
go install
```

Make sure your `~/go/bin` directory is in your system's PATH.

## Usage Examples

Here are the MVP commands supported:

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
Are you sure you want to delete all merged local branches? (y/n)
y
Cleaning branches...
> git branch --merged | grep -v "\*" | grep -v "master" | grep -v "main" | xargs -n 1 git branch -d
```

**4. Running tests and building**
```bash
$ commandspeak "run my tests and build"
Running tests and building...
> npm test && npm run build
```

## How to Extend

The codebase is structured to be easily extensible for hackathon teammates:
- **`internal/parser/parser.go`**: Add new regex patterns here to detect new intents and extract parameters.
- **`internal/executor/executor.go`**: Add new `switch` cases to map parsed intents to actual shell commands using `os/exec`.
- **`internal/config/config.go`**: Add new user preferences to the `Config` struct.
- **`main.go`**: The entry point for the CLI using `cobra`.
