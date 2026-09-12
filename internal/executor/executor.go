package executor

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"commandspeak/internal/config"
	"commandspeak/internal/parser"

	"github.com/fatih/color"
)

// ExecuteIntent takes a parsed intent and runs the appropriate shell commands.
func ExecuteIntent(intent parser.ParsedIntent, cfg *config.Config, globalDryRun bool) error {
	var intentName string
	isDryRun := intent.DryRun || globalDryRun

	switch intent.Type {
	case parser.IntentDeploy:
		intentName = "DEPLOY"
		if !isDryRun && !confirmGitClean() {
			return nil
		}
	case parser.IntentPush:
		intentName = "PUSH"
		if !isDryRun && !confirmGitClean() {
			return nil
		}
	case parser.IntentCleanBranches:
		intentName = "CLEAN_BRANCHES"
		if !isDryRun {
			color.Yellow("Are you sure you want to delete all merged branches? (y/n)")
			if !getUserConfirmation() {
				color.Yellow("Cancelled.")
				return nil
			}
			color.Cyan("Cleaning branches...")
		}
	case parser.IntentRunTests:
		intentName = "RUN_TESTS"
	case parser.IntentStatus:
		intentName = "STATUS"
	case parser.IntentInit:
		intentName = "INIT"
	case parser.IntentInstall:
		intentName = "INSTALL"
	case parser.IntentStart:
		intentName = "START"
	case parser.IntentClone:
		intentName = "CLONE"
		if intent.Parameters["repoUrl"] == "" && !isDryRun {
			fmt.Print("Enter the GitHub repo URL to clone: ")
			reader := bufio.NewReader(os.Stdin)
			url, _ := reader.ReadString('\n')
			url = strings.TrimSpace(url)
			intent.Parameters["repoUrl"] = url
		}
	case parser.IntentUndo:
		intentName = "UNDO"
		if !isDryRun {
			color.Yellow("Are you sure you want to undo your last commit? (y/n)")
			if !getUserConfirmation() {
				color.Yellow("Cancelled.")
				return nil
			}
		}
	default:
		return fmt.Errorf("unknown intent")
	}

	cmdTemplate := cfg.Commands[intentName]
	if cmdTemplate == "" {
		return fmt.Errorf("no command template configured for intent: %s", intentName)
	}

	// Replace placeholders like {{message}}, {{platform}}, {{repoUrl}}
	cmdStr := cmdTemplate
	for key, val := range intent.Parameters {
		cmdStr = strings.ReplaceAll(cmdStr, fmt.Sprintf("{{%s}}", key), val)
	}

	if isDryRun {
		color.HiBlack("Dry-run — would execute:")
		color.HiBlack("  > %s", cmdStr)
		return nil
	}

	return shellRun(".", cmdStr)
}

// shellRun runs a command string using the correct shell for the current OS.
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

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("command failed: %w", err)
	}
	return nil
}

// confirmGitClean checks `git status` for uncommitted changes and warns the user.
func confirmGitClean() bool {
	var out []byte
	var err error
	if runtime.GOOS == "windows" {
		out, err = exec.Command("cmd", "/c", "git status -s").Output()
	} else {
		out, err = exec.Command("bash", "-c", "git status -s").Output()
	}
	if err != nil {
		return true // not a git repo yet, let the command handle it
	}
	if len(strings.TrimSpace(string(out))) > 0 {
		color.Yellow("You have uncommitted changes:")
		fmt.Print(string(out))
		color.Yellow("Proceed anyway? (y/n)")
		return getUserConfirmation()
	}
	return true
}

func getUserConfirmation() bool {
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}
