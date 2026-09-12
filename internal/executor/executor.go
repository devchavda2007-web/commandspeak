package executor

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"commandspeak/internal/config"
	"commandspeak/internal/parser"
)

// ExecuteIntent takes a parsed intent and runs the appropriate shell commands using config templates.
func ExecuteIntent(intent parser.ParsedIntent, cfg *config.Config) error {
	var cmdTemplate string
	var intentName string

	switch intent.Type {
	case parser.IntentDeploy:
		intentName = "DEPLOY"
	case parser.IntentPush:
		intentName = "PUSH"
	case parser.IntentCleanBranches:
		intentName = "CLEAN_BRANCHES"
	case parser.IntentRunTests:
		intentName = "RUN_TESTS"
	default:
		return fmt.Errorf("unknown intent. Please try phrasing your command differently")
	}

	cmdTemplate = cfg.Commands[intentName]
	if cmdTemplate == "" {
		return fmt.Errorf("no command template configured for intent: %s", intentName)
	}

	// Handle confirmations
	if intent.Type == parser.IntentCleanBranches {
		fmt.Println("Are you sure you want to clean branches? (y/n)")
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		response = strings.ToLower(strings.TrimSpace(response))
		if response != "y" && response != "yes" {
			fmt.Println("Operation cancelled.")
			return nil
		}
		fmt.Println("Cleaning branches...")
	}

	// Replace placeholders like {{message}} and {{platform}}
	cmdStr := cmdTemplate
	for key, val := range intent.Parameters {
		placeholder := fmt.Sprintf("{{%s}}", key)
		cmdStr = strings.ReplaceAll(cmdStr, placeholder, val)
	}

	return runCommand(cmdStr)
}

// runCommand runs a command string using the system shell and prints output to stdout/stderr.
func runCommand(cmdStr string) error {
	fmt.Printf("> %s\n", cmdStr)
	
	cmd := exec.Command("bash", "-c", cmdStr)
	
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("command execution failed: %w", err)
	}

	return nil
}
