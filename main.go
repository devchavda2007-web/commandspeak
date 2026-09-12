package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"commandspeak/internal/config"
	"commandspeak/internal/executor"
	"commandspeak/internal/parser"

	"github.com/spf13/cobra"
)

var cfg *config.Config

func main() {
	// Initialize config
	var err error
	cfg, err = config.LoadConfig()
	if err != nil {
		fmt.Printf("Warning: Failed to load config: %v\n", err)
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
	fmt.Println("--- CommandSpeak Configuration ---")
	
	// List current commands
	for intent, cmdStr := range cfg.Commands {
		fmt.Printf("- %s:\n    %s\n", intent, cmdStr)
	}
	fmt.Println("----------------------------------")

	fmt.Print("\nWhich command do you want to change? (DEPLOY, PUSH, CLEAN_BRANCHES, RUN_TESTS) or type 'exit': ")
	intentToChange, _ := reader.ReadString('\n')
	intentToChange = strings.ToUpper(strings.TrimSpace(intentToChange))

	if intentToChange == "EXIT" || intentToChange == "" {
		fmt.Println("Cancelled.")
		return
	}

	_, exists := cfg.Commands[intentToChange]
	if !exists {
		fmt.Printf("Unknown intent '%s'.\n", intentToChange)
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
			fmt.Printf("Error saving config: %v\n", err)
		} else {
			fmt.Println("Configuration updated successfully!")
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
			fmt.Printf("Error reading input: %v\n", err)
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
		fmt.Println("Sorry, I couldn't understand that command. Try 'commandspeak --help' for examples.")
		return
	}

	err := executor.ExecuteIntent(intent, cfg)
	if err != nil {
		fmt.Printf("Error executing command: %v\n", err)
	} else {
		fmt.Println("Success!")
	}
}
