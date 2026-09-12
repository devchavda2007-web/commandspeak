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
	IntentUnknown       IntentType = "UNKNOWN"
)

type ParsedIntent struct {
	Type       IntentType
	Parameters map[string]string
}

// ParseSentence takes a natural language string and returns the parsed intent and parameters.
func ParseSentence(sentence string) ParsedIntent {
	sentence = strings.ToLower(strings.TrimSpace(sentence))

	// 1. Check for deploy intent
	// e.g., "deploy my project to vercel"
	deployRegex := regexp.MustCompile(`deploy(?:\s+my)?(?:\s+project)?(?:\s+to\s+(\w+))?`)
	if deployMatch := deployRegex.FindStringSubmatch(sentence); deployMatch != nil {
		platform := "vercel" // default
		if len(deployMatch) > 1 && deployMatch[1] != "" {
			platform = deployMatch[1]
		}
		return ParsedIntent{
			Type: IntentDeploy,
			Parameters: map[string]string{
				"platform": platform,
			},
		}
	}

	// 2. Check for push intent
	// e.g., "push my code with message initial commit"
	// Also handle quotes "push my code with message 'initial commit'"
	pushRegex := regexp.MustCompile(`push(?:\s+my)?(?:\s+code)?(?:\s+with\s+message\s+(.+))?`)
	if pushMatch := pushRegex.FindStringSubmatch(sentence); pushMatch != nil {
		message := "update" // default
		if len(pushMatch) > 1 && pushMatch[1] != "" {
			message = strings.Trim(pushMatch[1], "\"'")
		}
		return ParsedIntent{
			Type: IntentPush,
			Parameters: map[string]string{
				"message": message,
			},
		}
	}

	// 3. Check for clean branches intent
	// e.g., "clean my old branches"
	if strings.Contains(sentence, "clean") && strings.Contains(sentence, "branch") {
		return ParsedIntent{
			Type:       IntentCleanBranches,
			Parameters: make(map[string]string),
		}
	}

	// 4. Check for run tests intent
	// e.g., "run my tests and build"
	if strings.Contains(sentence, "test") && strings.Contains(sentence, "build") {
		return ParsedIntent{
			Type:       IntentRunTests,
			Parameters: make(map[string]string),
		}
	}

	return ParsedIntent{
		Type:       IntentUnknown,
		Parameters: make(map[string]string),
	}
}
