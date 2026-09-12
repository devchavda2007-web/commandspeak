package main

import (
	"bufio"
	"fmt"
	"os"
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
  commandspeak "push my code with message initial commit"`,
		Args: cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				// Interactive mode
				runInteractive()
				return
			}

			// Single command mode
			sentence := strings.Join(args, " ")
			processSentence(sentence)
		},
	}

	rootCmd.PersistentFlags().BoolVar(&globalDryRun, "dry-run", false, "Print the shell command without executing it")

	var changeCmd = &cobra.Command{
		Use:   "change",
		Short: "Interactively change the configured shell commands",
		Run: func(cmd *cobra.Command, args []string) {
			runChangeCommand()
		},
	}

	rootCmd.AddCommand(changeCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runChangeCommand() {
	reader := bufio.NewReader(os.Stdin)
	color.Cyan("--- CommandSpeak Configuration ---")
	
	// List current commands
	for intent, cmdStr := range cfg.Commands {
		fmt.Printf("- %s:\n    %s\n", intent, cmdStr)
	}
	fmt.Println("----------------------------------")

	fmt.Print("\nWhich command do you want to change? (DEPLOY, PUSH, CLEAN_BRANCHES, RUN_TESTS, STATUS, UNDO, INIT, INSTALL, START) or type 'exit': ")
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
	fmt.Print("\nEnter new command template (use {{message}} or {{platform}} where applicable):\n> ")
	
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

func runInteractive() {
	fmt.Println("Welcome to CommandSpeak! Type a natural language command (or 'exit' to quit):")
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
		if sentence == "" {
			continue
		}

		processSentence(sentence)
	}
}

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
