package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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
	var err error
	cfg, err = config.LoadConfig()
	if err != nil {
		color.Yellow("Warning: Failed to load config: %v", err)
	}

	// ─── Root command ─────────────────────────────────────────────────────────
	var rootCmd = &cobra.Command{
		Use:   "commandspeak [sentence]",
		Short: "CommandSpeak - Natural Language CLI for Developers",
		Long: `
  ██████╗ ███████╗██╗   ██╗
  ██╔══██╗██╔════╝██║   ██║
  ██║  ██║█████╗  ██║   ██║
  ██║  ██║██╔══╝  ╚██╗ ██╔╝
  ██████╔╝███████╗ ╚████╔╝
  ╚═════╝ ╚══════╝  ╚═══╝   CommandSpeak

Speak plain English. Run real developer commands.

  commandspeak "push my code with message fixed bug"
  commandspeak "deploy to vercel"
  commandspeak repo                          ← open GitHub repo manager
  commandspeak repo https://github.com/u/p  ← clone & manage a specific repo
  commandspeak change                        ← customize any command`,
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
	rootCmd.PersistentFlags().BoolVar(&globalDryRun, "dry-run", false, "Show command without executing it")

	// ─── change subcommand ────────────────────────────────────────────────────
	var changeCmd = &cobra.Command{
		Use:   "change",
		Short: "Customize what shell command runs for each intent",
		Run: func(cmd *cobra.Command, args []string) {
			runChangeCommand()
		},
	}

	// ─── repo subcommand ──────────────────────────────────────────────────────
	var repoURLFlag string
	var repoCmd = &cobra.Command{
		Use:   "repo [github-url]",
		Short: "Clone and manage any GitHub repository interactively",
		Long: `
CommandSpeak Repo Manager — Work with ANY GitHub repository.

Provide the repo URL in one of these ways:

  1. Direct argument:
       commandspeak repo https://github.com/user/project

  2. As a flag:
       commandspeak repo --url https://github.com/user/project
       commandspeak repo -u https://github.com/user/project

  3. Interactive (just run and paste when asked):
       commandspeak repo

  4. Inside interactive mode shell:
       > repo https://github.com/user/project

Once you provide a URL, you get a full menu to:
  ● Clone the repo to your machine
  ● Browse all files and folders
  ● View any file's contents
  ● Edit files in Notepad (Windows) or nano (Linux/Mac)
  ● Commit + push your changes back
  ● Pull latest changes
  ● Check git status
  ● Switch to a completely different repo URL`,
		Args: cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			urlFromArg := ""
			if len(args) == 1 {
				urlFromArg = args[0]
			}
			// Argument takes priority over flag
			if urlFromArg != "" {
				repoURLFlag = urlFromArg
			}
			runRepoManager(repoURLFlag)
		},
	}
	repoCmd.Flags().StringVarP(&repoURLFlag, "url", "u", "", "GitHub repository URL to clone/manage")

	// ─── ui subcommand ────────────────────────────────────────────────────────
	var uiCmd = &cobra.Command{
		Use:   "ui",
		Short: "Start the CommandSpeak local UI (history & settings)",
		Run: func(cmd *cobra.Command, args []string) {
			runUIServer()
		},
	}

	rootCmd.AddCommand(changeCmd)
	rootCmd.AddCommand(repoCmd)
	rootCmd.AddCommand(uiCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// =============================================================================
// REPO MANAGER
// =============================================================================

func runRepoManager(initialURL string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println()
	color.Cyan("╔══════════════════════════════════════════════════╗")
	color.Cyan("║       CommandSpeak — GitHub Repo Manager         ║")
	color.Cyan("╠══════════════════════════════════════════════════╣")
	color.Cyan("║  Clone • Browse • Edit • Commit • Push • Pull    ║")
	color.Cyan("╚══════════════════════════════════════════════════╝")
	fmt.Println()

	// ── STEP 1: Get Repo URL ──────────────────────────────────────────────────
	repoURL := strings.TrimSpace(initialURL)

	if repoURL == "" {
		color.Yellow("┌─ Enter the GitHub Repo URL to clone / work with:")
		color.HiBlack("│  Examples:")
		color.HiBlack("│    https://github.com/torvalds/linux")
		color.HiBlack("│    https://github.com/devchavda2007-web/commandspeak")
		color.HiBlack("│    https://github.com/anyuser/anyrepo.git")
		fmt.Print("└─ URL > ")
		input, _ := reader.ReadString('\n')
		repoURL = strings.TrimSpace(input)
	} else {
		color.Green("✔  Using repo: %s", repoURL)
	}

	if repoURL == "" {
		color.Red("✗  No URL entered. Exiting repo manager.")
		return
	}

	// ── STEP 2: Derive local folder name ─────────────────────────────────────
	repoName := strings.TrimSuffix(filepath.Base(repoURL), ".git")
	cloneDir, _ := filepath.Abs(repoName)

	color.HiBlack("  Local folder: %s", cloneDir)
	fmt.Println()

	// ── STEP 3: Clone if needed ───────────────────────────────────────────────
	if _, err := os.Stat(cloneDir); os.IsNotExist(err) {
		color.Cyan("⬇  Cloning repository, please wait...")
		if err := shellRun(".", fmt.Sprintf("git clone %s", repoURL)); err != nil {
			color.Red("✗  Clone failed: %v", err)
			color.HiBlack("   Make sure the URL is correct and git is installed.")
			return
		}
		color.Green("✔  Cloned successfully into: %s", cloneDir)
	} else {
		color.Green("✔  Repo already exists at: %s  (skipping clone)", cloneDir)
	}

	// ── STEP 4: Interactive Menu ──────────────────────────────────────────────
	for {
		fmt.Println()
		color.Cyan("╔══════════════════════════════════════════════════╗")
		fmt.Printf("  Repo: %s\n", color.CyanString(repoURL))
		fmt.Printf("  Dir:  %s\n", color.HiBlackString(cloneDir))
		color.Cyan("╠══════════════════════════════════════════════════╣")
		fmt.Println("  1.  📂  List all files & folders")
		fmt.Println("  2.  📄  View a file")
		fmt.Println("  3.  ✏️   Edit a file")
		fmt.Println("  4.  🚀  Commit & push changes")
		fmt.Println("  5.  ⬇️   Pull latest from remote")
		fmt.Println("  6.  🔍  Show git status")
		fmt.Println("  7.  📝  Show git log (last 10 commits)")
		fmt.Println("  8.  🌿  Show all branches")
		fmt.Println("  9.  🔀  Switch to a different repo URL")
		fmt.Println("  10. 🗑️   Delete a file from the repo")
		fmt.Println("  0.  ←   Exit repo manager")
		color.Cyan("╚══════════════════════════════════════════════════╝")
		fmt.Print("  Choice > ")

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
			color.Cyan("⬇  Pulling latest changes...")
			shellRun(cloneDir, "git pull")
		case "6":
			shellRun(cloneDir, "git status")
		case "7":
			shellRun(cloneDir, "git log --oneline -10")
		case "8":
			shellRun(cloneDir, "git branch -a")
		case "9":
			repoURL, cloneDir, repoName = repoSwitchURL(reader, repoURL)
		case "10":
			repoDeleteFile(cloneDir, reader)
		case "0":
			color.Green("✔  Exited repo manager. Your files are at: %s", cloneDir)
			return
		default:
			color.Red("✗  Unknown option '%s'. Enter a number from 0-10.", choice)
		}
	}
}

// ─── Option 1: List files ─────────────────────────────────────────────────────
func repoListFiles(dir string) {
	fmt.Println()
	color.Cyan("📂  Files in repository:")
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

// ─── Option 2: View file ──────────────────────────────────────────────────────
func repoViewFile(dir string, reader *bufio.Reader) {
	fmt.Print("  📄 File path (relative to repo root) > ")
	filePath, _ := reader.ReadString('\n')
	filePath = strings.TrimSpace(filePath)

	fullPath := filepath.Join(dir, filePath)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		color.Red("✗  Could not read file '%s': %v", filePath, err)
		return
	}
	fmt.Println()
	color.Cyan("─── %s ───────────────────────────────────────────────", filePath)
	fmt.Println(string(data))
	color.Cyan("────────────────────────────────────────────────────────")
}

// ─── Option 3: Edit file ──────────────────────────────────────────────────────
func repoEditFile(dir string, reader *bufio.Reader) {
	fmt.Print("  ✏️  File path to edit (relative to repo root) > ")
	filePath, _ := reader.ReadString('\n')
	filePath = strings.TrimSpace(filePath)
	fullPath := filepath.Join(dir, filePath)

	// Make sure the file exists before opening
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		color.Yellow("File does not exist. Create it? (y/n)")
		ans, _ := reader.ReadString('\n')
		ans = strings.TrimSpace(strings.ToLower(ans))
		if ans != "y" && ans != "yes" {
			return
		}
		os.MkdirAll(filepath.Dir(fullPath), 0755)
		os.WriteFile(fullPath, []byte(""), 0644)
	}

	color.Cyan("✏️  Opening %s ...", fullPath)

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		// On Windows: open with notepad natively (no bash needed)
		cmd = exec.Command("notepad", fullPath)
	} else {
		// On Linux/Mac: use nano
		cmd = exec.Command("nano", fullPath)
		cmd.Stdin = os.Stdin
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		color.Red("✗  Could not open editor: %v", err)
		color.HiBlack("   File is at: %s  — edit it manually.", fullPath)
	}
}

// ─── Option 4: Commit & Push ──────────────────────────────────────────────────
func repoCommitAndPush(dir string, reader *bufio.Reader) {
	fmt.Print("  🚀 Commit message > ")
	msg, _ := reader.ReadString('\n')
	msg = strings.TrimSpace(msg)
	if msg == "" {
		msg = "update via CommandSpeak"
	}

	color.Cyan("  Staging all changes...")
	if err := shellRun(dir, "git add ."); err != nil {
		color.Red("✗  git add failed: %v", err)
		return
	}

	color.Cyan("  Committing...")
	commitCmd := fmt.Sprintf(`git commit -m "%s"`, msg)
	if err := shellRun(dir, commitCmd); err != nil {
		color.Yellow("  Nothing new to commit, or commit failed.")
		return
	}

	color.Cyan("  Pushing...")
	if err := shellRun(dir, "git push"); err != nil {
		color.Red("✗  Push failed. You may not have write access to this repo.")
		color.HiBlack("   If this is someone else's repo, fork it on GitHub first.")
	} else {
		color.Green("✔  Changes pushed successfully!")
	}
}

// ─── Option 9: Switch repo URL ────────────────────────────────────────────────
func repoSwitchURL(reader *bufio.Reader, currentURL string) (newURL string, newDir string, newName string) {
	color.Yellow("  Current repo: %s", currentURL)
	fmt.Print("  🔀 Enter new GitHub repo URL > ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		color.Red("No URL entered, keeping current repo.")
		return currentURL, ".", filepath.Base(currentURL)
	}
	name := strings.TrimSuffix(filepath.Base(input), ".git")
	dir, _ := filepath.Abs(name)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		color.Cyan("⬇  Cloning %s ...", input)
		if err := shellRun(".", fmt.Sprintf("git clone %s", input)); err != nil {
			color.Red("✗  Clone failed: %v", err)
			return currentURL, ".", filepath.Base(currentURL)
		}
	}
	color.Green("✔  Switched to: %s", input)
	return input, dir, name
}

// ─── Option 10: Delete file ───────────────────────────────────────────────────
func repoDeleteFile(dir string, reader *bufio.Reader) {
	fmt.Print("  🗑️  File path to delete (relative to repo root) > ")
	filePath, _ := reader.ReadString('\n')
	filePath = strings.TrimSpace(filePath)
	fullPath := filepath.Join(dir, filePath)

	color.Yellow("  Are you sure you want to delete '%s'? (y/n)", filePath)
	ans, _ := reader.ReadString('\n')
	ans = strings.TrimSpace(strings.ToLower(ans))
	if ans != "y" && ans != "yes" {
		fmt.Println("  Cancelled.")
		return
	}
	if err := os.Remove(fullPath); err != nil {
		color.Red("✗  Could not delete: %v", err)
	} else {
		color.Green("✔  Deleted: %s", filePath)
	}
}

// =============================================================================
// SHELL HELPERS  (works on Windows AND Linux/Mac)
// =============================================================================

func shellRun(dir string, cmdStr string) error {
	color.HiBlack("  > %s", cmdStr)

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", cmdStr)
	} else {
		cmd = exec.Command("bash", "-c", cmdStr)
	}

	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// =============================================================================
// CHANGE COMMAND
// =============================================================================

func runChangeCommand() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println()
	color.Cyan("╔══════════════════════════════════════════════════╗")
	color.Cyan("║       CommandSpeak — Command Configuration       ║")
	color.Cyan("╚══════════════════════════════════════════════════╝")
	fmt.Println()
	color.HiBlack("  Current command templates:")
	fmt.Println()

	for intent, cmdStr := range cfg.Commands {
		color.Yellow("  %-16s", intent)
		color.HiBlack("    → %s", cmdStr)
	}

	fmt.Println()
	fmt.Print("  Which intent to change?\n  (DEPLOY / PUSH / CLEAN_BRANCHES / RUN_TESTS / STATUS / UNDO / INIT / INSTALL / START / CLONE)\n  > ")
	intentToChange, _ := reader.ReadString('\n')
	intentToChange = strings.ToUpper(strings.TrimSpace(intentToChange))

	if intentToChange == "EXIT" || intentToChange == "" {
		fmt.Println("  Cancelled.")
		return
	}

	_, exists := cfg.Commands[intentToChange]
	if !exists {
		color.Red("✗  Unknown intent '%s'.", intentToChange)
		return
	}

	fmt.Println()
	color.Yellow("  Current template for %s:", intentToChange)
	fmt.Printf("  %s\n", cfg.Commands[intentToChange])
	fmt.Println()
	color.HiBlack("  Placeholders you can use: {{message}}  {{platform}}  {{repoUrl}}")
	fmt.Print("  New command > ")

	newCmd, _ := reader.ReadString('\n')
	newCmd = strings.TrimSpace(newCmd)

	if newCmd != "" {
		cfg.Commands[intentToChange] = newCmd
		if err := config.SaveConfig(cfg); err != nil {
			color.Red("✗  Error saving config: %v", err)
		} else {
			color.Green("✔  Saved! '%s' will now run: %s", intentToChange, newCmd)
		}
	} else {
		fmt.Println("  No changes made.")
	}
}

// =============================================================================
// INTERACTIVE MODE
// =============================================================================

func runInteractive() {
	fmt.Println()
	color.Cyan("╔══════════════════════════════════════════════════╗")
	color.Cyan("║         CommandSpeak — Interactive Shell         ║")
	color.Cyan("╠══════════════════════════════════════════════════╣")
	color.HiBlack("║  Type a sentence or one of these shortcuts:      ║")
	color.HiBlack("║    repo                   → open repo manager    ║")
	color.HiBlack("║    repo <url>             → clone & manage repo  ║")
	color.HiBlack("║    exit / quit            → close                ║")
	color.Cyan("╚══════════════════════════════════════════════════╝")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		sentence, err := reader.ReadString('\n')
		if err != nil {
			color.Red("Error reading input: %v", err)
			continue
		}

		sentence = strings.TrimSpace(sentence)

		switch {
		case sentence == "":
			continue
		case sentence == "exit" || sentence == "quit":
			color.Green("Goodbye! 👋")
			return
		case sentence == "repo":
			runRepoManager("")
		case strings.HasPrefix(sentence, "repo "):
			inlineURL := strings.TrimSpace(strings.TrimPrefix(sentence, "repo "))
			runRepoManager(inlineURL)
		default:
			processSentence(sentence)
		}
	}
}

// =============================================================================
// PROCESS SENTENCE
// =============================================================================

func processSentence(sentence string) {
	intent := parser.ParseSentence(sentence)

	if intent.Type == parser.IntentUnknown {
		color.Red("✗  Sorry, I couldn't understand: \"%s\"", sentence)
		fmt.Println()
		color.Yellow("  Try one of these:")
		examples := []string{
			`"deploy my project to vercel"`,
			`"push my code with message fixed bug"`,
			`"clean my old branches"`,
			`"run my tests and build"`,
			`"show me my git status"`,
			`"undo my last commit"`,
			`"initialize git repository"`,
			`"install dependencies"`,
			`"start server"`,
			`"clone repo https://github.com/user/project"`,
			`"repo" or "repo <url>"  ← to open the full repo manager`,
		}
		for _, ex := range examples {
			fmt.Printf("    %s\n", ex)
		}
		return
	}

	err := executor.ExecuteIntent(intent, cfg, globalDryRun)
	if err != nil {
		color.Red("✗  Error: %v", err)
	} else {
		if !intent.DryRun && !globalDryRun {
			color.Green("✔  Done!")
		}
	}
}

// =============================================================================
// UI LAUNCHER
// =============================================================================

func runUIServer() {
	color.Cyan("\nStarting CommandSpeak Frontend...")
	color.HiBlack("Opening index.html directly...")

	absPath, err := filepath.Abs("index.html")
	if err != nil {
		color.Red("Failed to find index.html: %v", err)
		return
	}
	openBrowser("file:///" + filepath.ToSlash(absPath))
	color.Green("✔ Opened in browser successfully.")
}

func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		color.Red("Could not open browser automatically: %v", err)
		color.White("Please open: %s", url)
	}
}
