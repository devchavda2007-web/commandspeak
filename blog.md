# Building a Full GitHub Repo Manager in the CLI with CommandSpeak

When building CLI tools for developers, the goal is always to reduce friction. Switching contexts between the terminal, browser, and an IDE just to make a quick edit on a GitHub repository is a pain point we've all felt. 

To solve this, we've introduced a major feature to **CommandSpeak**: a fully interactive GitHub Repository Manager directly in your terminal.

## What's New?

We've added a new `repo` subcommand and a natural language `CLONE` intent that lets you download, browse, and edit repositories without leaving the command line.

### 1. Flexible URL Input
We wanted to make sure you could jump into a repository however you prefer:
- **Direct Argument**: `commandspeak repo https://github.com/user/project`
- **Flag**: `commandspeak repo -u https://github.com/user/project`
- **Interactive Prompt**: Run `commandspeak repo` and paste when asked.
- **Natural Language**: Just say `commandspeak "clone repo https://github.com/user/project"`

### 2. The 10-Option Interactive Menu
Once you provide a repository URL, CommandSpeak automatically clones it (or detects if it's already cloned) and drops you into a sleek interactive menu. You get 10 powerful options:
1. 📂 **List all files & folders**: Navigate the project structure.
2. 📄 **View a file**: Quickly read file contents.
3. ✏️ **Edit a file**: Opens Notepad (Windows) or Nano (Linux/Mac) for quick changes.
4. 🚀 **Commit & push changes**: Stages all changes, asks for a commit message, and pushes.
5. ⬇️ **Pull latest changes**: Sync with the remote repo.
6. 🔍 **Show git status**: See what's modified.
7. 📝 **Show git log**: View recent commits.
8. 🌿 **Show all branches**: See the branches.
9. 🔀 **Switch repo URL**: Jump to another project without restarting the CLI.
10. 🗑️ **Delete a file**: Remove files you don't need.

### 3. Cross-Platform & Windows Compatibility
We made sure all the underlying shell commands are OS-aware. Whether you're on a bash shell in Linux/Mac or running `cmd.exe` on Windows, CommandSpeak detects your environment and executes the correct native commands (like defaulting to `notepad` on Windows).

## What's Next?
This update transforms CommandSpeak from a simple command runner into a comprehensive workspace manager. During this 24-hour hackathon, our focus is on refining these developer workflows to make them as seamless as possible.

Try it out by running `commandspeak repo`!
