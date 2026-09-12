package parser

import (
	"testing"
)

func TestParseSentence(t *testing.T) {
	tests := []struct {
		name           string
		sentence       string
		expectedIntent IntentType
		expectedDryRun bool
		expectedParams map[string]string
	}{
		{
			name:           "Basic Deploy",
			sentence:       "deploy my project to vercel",
			expectedIntent: IntentDeploy,
			expectedDryRun: false,
			expectedParams: map[string]string{"platform": "vercel"},
		},
		{
			name:           "Synonym Deploy",
			sentence:       "deploy to netlify now",
			expectedIntent: IntentDeploy,
			expectedDryRun: false,
			expectedParams: map[string]string{"platform": "netlify"},
		},
		{
			name:           "Basic Push with parameter",
			sentence:       "push my code with message 'fixed bug'",
			expectedIntent: IntentPush,
			expectedDryRun: false,
			expectedParams: map[string]string{"message": "fixed bug"},
		},
		{
			name:           "Reordered Push",
			sentence:       "with message update push code",
			expectedIntent: IntentPush,
			expectedDryRun: false,
			expectedParams: map[string]string{"message": "update"},
		},
		{
			name:           "Dry Run detection",
			sentence:       "show me how to run my tests and build",
			expectedIntent: IntentRunTests,
			expectedDryRun: true,
			expectedParams: map[string]string{},
		},
		{
			name:           "Unknown Intent Fallback",
			sentence:       "make me a sandwich please",
			expectedIntent: IntentUnknown,
			expectedDryRun: false,
			expectedParams: map[string]string{},
		},
		{
			name:           "Status Intent",
			sentence:       "what is my git status",
			expectedIntent: IntentStatus,
			expectedDryRun: false,
			expectedParams: map[string]string{},
		},
		{
			name:           "Undo Intent",
			sentence:       "undo my last commit",
			expectedIntent: IntentUndo,
			expectedDryRun: false,
			expectedParams: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			intent := ParseSentence(tt.sentence)

			if intent.Type != tt.expectedIntent {
				t.Errorf("expected intent %s, got %s", tt.expectedIntent, intent.Type)
			}
			
			if intent.DryRun != tt.expectedDryRun {
				t.Errorf("expected dryRun %v, got %v", tt.expectedDryRun, intent.DryRun)
			}

			for k, v := range tt.expectedParams {
				if intent.Parameters[k] != v {
					t.Errorf("expected parameter %s to be %s, got %s", k, v, intent.Parameters[k])
				}
			}
		})
	}
}
