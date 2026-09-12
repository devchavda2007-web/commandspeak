package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"commandspeak/internal/config"
	"commandspeak/internal/executor"
	"commandspeak/internal/parser"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var cfg *config.Config
var globalDryRun bool

func main() {
	// Initialize config
	var err error
	cfg, err = config.LoadConfig()
	if err != nil {
		color.Yellow("Warning: Failed to load config: %v", err)
	}

	var rootCmd = &cobra.Command{
		Use:   "commandspeak [sentence]",
		Short: "CommandSpeak - Natural Language CLI for Developers",
		Long: `CommandSpeak allows you to run complex terminal commands using simple natural language sentences.
For example: 
  commandspeak "deploy my project to vercel"
  commandspeak "push my code with message initial commit"
  commandspeak repo   (open the GitHub repo manager)`,
		Args: cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				runInteractive()
				return
			}
			sentence := strings.Join(args, " ")
			processSentence(sentence)
		},
	}

	rootCmd.PersistentFlags().BoolVar(&globalDryRun, "dry-run", false, "Print the shell command without executing it")

	// --- change subcommand ---
	var changeCmd = &cobra.Command{
		Use:   "change",
		Short: "Interactively change the configured shell commands",
		Run: func(cmd *cobra.Command, args []string) {
			runChangeCommand()
		},
	}

	// --- repo subcommand ---
	var repoCmd = &cobra.Command{
		Use:   "repo",
		Short: "Manage any GitHub repository (clone, browse files, push changes)",
		Long: `The repo command lets you input any GitHub repo URL and interact with it:
  - Clone / download it to your machine
  - List its files
  - Open and view a file
  - Edit a file directly in the terminal
  - Commit and push your changes back`,
		Run: func(cmd *cobra.Command, args []string) {
			runRepoManager()
		},
	}

	rootCmd.AddCommand(changeCmd)
	rootCmd.AddCommand(repoCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// REPO MANAGER
// ─────────────────────────────────────────────────────────────────────────────

func runRepoManager() {
	reader := bufio.NewReader(os.Stdin)

	color.Cyan("╔══════════════════════════════════════════╗")
	color.Cyan("║     CommandSpeak — GitHub Repo Manager   ║")
	color.Cyan("╚══════════════════════════════════════════╝")
	fmt.Println()

	// Step 1: Get the repo URL
	color.Yellow("Enter the GitHub repo URL you want to work with:")
	fmt.Print("  URL: ")
	repoURL, _ := reader.ReadString('\n')
	repoURL = strings.TrimSpace(repoURL)

	if repoURL == "" {
		color.Red("No URL entered. Exiting.")
		return
	}

	// Derive a local folder name from the URL  e.g. https://github.com/user/project → project
	repoName := strings.TrimSuffix(filepath.Base(repoURL), ".git")
	cloneDir := filepath.Join(".", repoName)

	// Step 2: Clone if not already cloned
	if _, err := os.Stat(cloneDir); os.IsNotExist(err) {
		color.Cyan("\nCloning repository...")
		if err := runShell(fmt.Sprintf("git clone %s", repoURL)); err != nil {
			color.Red("Failed to clone: %v", err)
			return
		}
		color.Green("Cloned into ./%s", repoName)
	} else {
		color.Green("Folder ./%s already exists — skipping clone.", repoName)
	}

	// Step 3: Interactive menu
	for {
		fmt.Println()
		color.Cyan("─── What do you want to do? ───────────────────────")
		fmt.Println("  1. List files in the repository")
		fmt.Println("  2. View a file")
		fmt.Println("  3. Edit a file (opens nano / notepad)")
		fmt.Println("  4. Commit and push changes")
		fmt.Println("  5. Pull latest changes from remote")
		fmt.Println("  6. Show git status")
		fmt.Println("  7. Change the remote URL (point to a different repo)")
		fmt.Println("  8. Exit repo manager")
		color.Cyan("───────────────────────────────────────────────────")
		fmt.Print("  Choice: ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			repoListFiles(cloneDir)
		case "2":
			repoViewFile(cloneDir, reader)
		case "3":
			repoEditFile(cloneDir, reader)
		case "4":
			repoCommitAndPush(cloneDir, reader)
		case "5":
			color.Cyan("Pulling latest changes...")
			runShellIn(cloneDir, "git pull")
		case "6":
			runShellIn(cloneDir, "git status")
		case "7":
			repoChangeRemote(cloneDir, reader)
		case "8":
			color.Green("Exited repo manager. Files saved in ./%s", repoName)
			return
		default:
			color.Red("Unknown option '%s'. Please enter 1-8.", choice)
		}
	}
}

func repoListFiles(dir string) {
	color.Cyan("\nFiles in repository:")
	entries, err := listDirRecursive(dir, 0)
	if err != nil {
		color.Red("Error listing files: %v", err)
		return
	}
	for _, e := range entries {
		fmt.Println(" ", e)
	}
}

func listDirRecursive(root string, depth int) ([]string, error) {
	var results []string
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		// Skip .git directory
		if e.Name() == ".git" {
			continue
		}
		indent := strings.Repeat("  ", depth)
		if e.IsDir() {
			results = append(results, fmt.Sprintf("%s📁 %s/", indent, e.Name()))
			sub, _ := listDirRecursive(filepath.Join(root, e.Name()), depth+1)
			results = append(results, sub...)
		} else {
			results = append(results, fmt.Sprintf("%s📄 %s", indent, e.Name()))
		}
	}
	return results, nil
}

func repoViewFile(dir string, reader *bufio.Reader) {
	fmt.Print("  Enter file path to view (relative to repo root): ")
	filePath, _ := reader.ReadString('\n')
	filePath = strings.TrimSpace(filePath)

	fullPath := filepath.Join(dir, filePath)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		color.Red("Could not read file: %v", err)
		return
	}
	color.Cyan("\n─── %s ─────────────────────────────────────────", filePath)
	fmt.Println(string(data))
	color.Cyan("──────────────────────────────────────────────────")
}

func repoEditFile(dir string, reader *bufio.Reader) {
	fmt.Print("  Enter file path to edit (relative to repo root): ")
	filePath, _ := reader.ReadString('\n')
	filePath = strings.TrimSpace(filePath)

	fullPath := filepath.Join(dir, filePath)

	// Try to open with notepad on Windows, nano on Linux/Mac
	var editorCmd string
	if _, err := exec.LookPath("nano"); err == nil {
		editorCmd = fmt.Sprintf("nano %s", fullPath)
	} else if _, err := exec.LookPath("notepad"); err == nil {
		editorCmd = fmt.Sprintf("notepad %s", fullPath)
	} else {
		color.Red("No editor found (tried nano and notepad). Please edit the file manually at: %s", fullPath)
		return
	}

	color.Cyan("Opening %s in editor...", fullPath)
	cmd := exec.Command("bash", "-c", editorCmd)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		// Try native Windows open for notepad
		winCmd := exec.Command("cmd", "/c", "notepad", fullPath)
		winCmd.Stdin = os.Stdin
		winCmd.Stdout = os.Stdout
		winCmd.Stderr = os.Stderr
		winCmd.Run()
	}
}

func repoCommitAndPush(dir string, reader *bufio.Reader) {
	fmt.Print("  Commit message: ")
	msg, _ := reader.ReadString('\n')
	msg = strings.TrimSpace(msg)
	if msg == "" {
		msg = "update via CommandSpeak"
	}

	color.Cyan("Committing and pushing...")
	cmd := fmt.Sprintf("git add . && git commit -m \"%s\" && git push", msg)
	if err := runShellIn(dir, cmd); err != nil {
		color.Red("Push failed: %v", err)
	} else {
		color.Green("Changes pushed successfully!")
	}
}

func repoChangeRemote(dir string, reader *bufio.Reader) {
	fmt.Print("  Enter new remote URL: ")
	newURL, _ := reader.ReadString('\n')
	newURL = strings.TrimSpace(newURL)
	if newURL == "" {
		color.Red("No URL entered.")
		return
	}
	cmd := fmt.Sprintf("git remote set-url origin %s", newURL)
	if err := runShellIn(dir, cmd); err != nil {
		color.Red("Failed to change remote: %v", err)
	} else {
		color.Green("Remote URL updated to: %s", newURL)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// SHELL HELPERS
// ─────────────────────────────────────────────────────────────────────────────

// runShell runs a bash command in the current directory.
func runShell(cmdStr string) error {
	color.HiBlack("> %s", cmdStr)
	cmd := exec.Command("bash", "-c", cmdStr)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// runShellIn runs a bash command inside a specific directory.
func runShellIn(dir string, cmdStr string) error {
	color.HiBlack("> (in %s) %s", dir, cmdStr)
	cmd := exec.Command("bash", "-c", cmdStr)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// ─────────────────────────────────────────────────────────────────────────────
// CHANGE COMMAND
// ─────────────────────────────────────────────────────────────────────────────

func runChangeCommand() {
	reader := bufio.NewReader(os.Stdin)
	color.Cyan("--- CommandSpeak Configuration ---")

	// List current commands
	for intent, cmdStr := range cfg.Commands {
		fmt.Printf("- %s:\n    %s\n", intent, cmdStr)
	}
	fmt.Println("----------------------------------")

	fmt.Print("\nWhich command do you want to change? (DEPLOY, PUSH, CLEAN_BRANCHES, RUN_TESTS, STATUS, UNDO, INIT, INSTALL, START, CLONE) or type 'exit': ")
	intentToChange, _ := reader.ReadString('\n')
	intentToChange = strings.ToUpper(strings.TrimSpace(intentToChange))

	if intentToChange == "EXIT" || intentToChange == "" {
		fmt.Println("Cancelled.")
		return
	}

	_, exists := cfg.Commands[intentToChange]
	if !exists {
		color.Red("Unknown intent '%s'.", intentToChange)
		return
	}

	fmt.Printf("\nCurrent command for %s is:\n%s\n", intentToChange, cfg.Commands[intentToChange])
	fmt.Print("\nEnter new command template (use {{message}}, {{platform}}, {{repoUrl}} where applicable):\n> ")

	newCmd, _ := reader.ReadString('\n')
	newCmd = strings.TrimSpace(newCmd)

	if newCmd != "" {
		cfg.Commands[intentToChange] = newCmd
		err := config.SaveConfig(cfg)
		if err != nil {
			color.Red("Error saving config: %v", err)
		} else {
			color.Green("Configuration updated successfully!")
		}
	} else {
		fmt.Println("No changes made.")
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// INTERACTIVE MODE
// ─────────────────────────────────────────────────────────────────────────────

func runInteractive() {
	color.Cyan("Welcome to CommandSpeak! Type a sentence or 'repo' to open the repo manager (type 'exit' to quit):")
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("\n> ")
		sentence, err := reader.ReadString('\n')
		if err != nil {
			color.Red("Error reading input: %v", err)
			continue
		}

		sentence = strings.TrimSpace(sentence)
		if sentence == "exit" || sentence == "quit" {
			fmt.Println("Goodbye!")
			break
		}
		if sentence == "repo" {
			runRepoManager()
			continue
		}
		if sentence == "" {
			continue
		}

		processSentence(sentence)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// PROCESS SENTENCE
// ─────────────────────────────────────────────────────────────────────────────

func processSentence(sentence string) {
	intent := parser.ParseSentence(sentence)

	if intent.Type == parser.IntentUnknown {
		color.Red("Sorry, I couldn't understand that command.")
		fmt.Println("\nHere are some supported examples you can try:")
		fmt.Println("  - \"deploy my project to vercel\"")
		fmt.Println("  - \"push my code with message update\"")
		fmt.Println("  - \"clean my old branches\"")
		fmt.Println("  - \"run my tests and build\"")
		fmt.Println("  - \"show me my git status\"")
		fmt.Println("  - \"undo my last commit\"")
		fmt.Println("  - \"initialize git repository\"")
		fmt.Println("  - \"install dependencies\"")
		fmt.Println("  - \"start server\"")
		fmt.Println("  - \"clone repo https://github.com/user/project\"")
		fmt.Println("  - run 'commandspeak repo' to open the full repo manager")
		return
	}

	err := executor.ExecuteIntent(intent, cfg, globalDryRun)
	if err != nil {
		color.Red("Error executing command: %v", err)
	} else {
		if !intent.DryRun && !globalDryRun {
			color.Green("Success!")
		}
	}
}
