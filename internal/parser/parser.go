package parser

import (
	"regexp"
	"strings"
)

type IntentType string

const (
	IntentDeploy        IntentType = "DEPLOY"
	IntentPush          IntentType = "PUSH"
	IntentCleanBranches IntentType = "CLEAN_BRANCHES"
	IntentRunTests      IntentType = "RUN_TESTS"
	IntentStatus        IntentType = "STATUS"
	IntentUndo          IntentType = "UNDO"
	IntentInit          IntentType = "INIT"
	IntentInstall       IntentType = "INSTALL"
	IntentStart         IntentType = "START"
	IntentClone         IntentType = "CLONE"
	IntentUnknown       IntentType = "UNKNOWN"
)

type ParsedIntent struct {
	Type       IntentType
	Parameters map[string]string
	DryRun     bool
}

// ParseSentence takes a natural language string and returns the parsed intent and parameters using flexible keyword scoring.
func ParseSentence(sentence string) ParsedIntent {
	sentence = strings.ToLower(strings.TrimSpace(sentence))
	
	intent := ParsedIntent{
		Type:       IntentUnknown,
		Parameters: make(map[string]string),
		DryRun:     false,
	}

	// 1. Check for Dry Run phrases
	if strings.HasPrefix(sentence, "show me") || strings.Contains(sentence, "show me") {
		intent.DryRun = true
		sentence = strings.ReplaceAll(sentence, "show me", "") // remove so it doesn't interfere with scoring
		sentence = strings.TrimSpace(sentence)
	}

	// 2. Keyword Scoring
	scores := make(map[IntentType]int)
	
	// Define keyword sets and their weight (1 point per word matched)
	keywords := map[IntentType][]string{
		IntentDeploy:        {"deploy", "vercel"},
		IntentPush:          {"push", "code", "message"},
		IntentCleanBranches: {"clean", "branch", "branches", "old"},
		IntentRunTests:      {"run", "test", "tests", "build"},
		IntentStatus:        {"status", "git"},
		IntentUndo:          {"undo", "commit", "last"},
		IntentInit:          {"initialize", "init", "repository"},
		IntentInstall:       {"install", "dependencies", "packages"},
		IntentStart:         {"start", "server", "dev"},
		IntentClone:         {"clone", "download", "fork"},
	}

	words := strings.Fields(sentence)
	for _, word := range words {
		// simple stemming/cleaning could happen here, but direct match is fine for MVP
		for intentType, intentKeywords := range keywords {
			for _, kw := range intentKeywords {
				if word == kw {
					scores[intentType]++
				}
			}
		}
	}

	// Find the intent with the highest score
	var bestIntent IntentType = IntentUnknown
	maxScore := 0
	for intentType, score := range scores {
		if score > maxScore {
			maxScore = score
			bestIntent = intentType
		}
	}

	// Require at least a score of 1 to match an intent
	if maxScore > 0 {
		intent.Type = bestIntent
	}

	// 3. Extract Parameters using Regex if intent requires it
	if intent.Type == IntentDeploy {
		intent.Parameters["platform"] = "vercel" // default
		deployRegex := regexp.MustCompile(`to\s+(\w+)`)
		if match := deployRegex.FindStringSubmatch(sentence); match != nil {
			intent.Parameters["platform"] = match[1]
		}
	} else if intent.Type == IntentPush {
		intent.Parameters["message"] = "update" // default
		pushRegex := regexp.MustCompile(`message\s+(.+)`)
		if match := pushRegex.FindStringSubmatch(sentence); match != nil {
			intent.Parameters["message"] = strings.Trim(match[1], "\"'")
		}
	} else if intent.Type == IntentClone {
		// Extract a GitHub URL from the sentence
		urlRegex := regexp.MustCompile(`(https?://\S+)`)
		if match := urlRegex.FindStringSubmatch(sentence); match != nil {
			intent.Parameters["repoUrl"] = match[1]
		}
	}

	return intent
}
