package executor

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"commandspeak/internal/config"
	"commandspeak/internal/parser"

	"github.com/fatih/color"
)

// ExecuteIntent takes a parsed intent and runs the appropriate shell commands using config templates.
// It also checks the dryRun flag and handles safety confirmations.
func ExecuteIntent(intent parser.ParsedIntent, cfg *config.Config, globalDryRun bool) error {
	var cmdTemplate string
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
			color.Yellow("Are you sure you want to clean branches? (y/n)")
			if !getUserConfirmation() {
				color.Yellow("Operation cancelled.")
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
			// Prompt user for the URL if not extracted from sentence
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
				color.Yellow("Operation cancelled.")
				return nil
			}
		}
	default:
		return fmt.Errorf("unknown intent")
	}

	cmdTemplate = cfg.Commands[intentName]
	if cmdTemplate == "" {
		return fmt.Errorf("no command template configured for intent: %s", intentName)
	}

	// Replace placeholders like {{message}} and {{platform}}
	cmdStr := cmdTemplate
	for key, val := range intent.Parameters {
		placeholder := fmt.Sprintf("{{%s}}", key)
		cmdStr = strings.ReplaceAll(cmdStr, placeholder, val)
	}

	if isDryRun {
		color.HiBlack("Dry-run mode. Would execute:")
		color.HiBlack("> %s\n", cmdStr)
		return nil
	}

	return runCommand(cmdStr)
}

// runCommand runs a command string using the system shell and prints output to stdout/stderr.
func runCommand(cmdStr string) error {
	color.HiBlack("> %s\n", cmdStr)
	
	cmd := exec.Command("bash", "-c", cmdStr)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("command execution failed: %w", err)
	}

	return nil
}

// confirmGitClean checks `git status -s` and prompts if there are uncommitted changes.
func confirmGitClean() bool {
	out, err := exec.Command("bash", "-c", "git status -s").Output()
	if err != nil {
		// Ignore error, might not be a git repo yet, let the underlying command fail if needed
		return true
	}
	
	if len(strings.TrimSpace(string(out))) > 0 {
		color.Yellow("You have uncommitted changes:")
		fmt.Print(string(out))
		color.Yellow("Do you want to proceed anyway? (y/n)")
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
