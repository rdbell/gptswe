package main

import (
	"errors"
	"strings"
	"testing"
)

func TestBuildPrompt(t *testing.T) {
	fileContents := "PRINT prompt.go\nFILE_START\n...\nFILE_END"
	tests := []struct {
		name                   string
		choice                 int
		detailsFlag            string
		expectedErr            error
		expectResponseContains string
	}{
		{
			name:                   "AskQuestion",
			choice:                 AskQuestion,
			detailsFlag:            "",
			expectedErr:            nil,
			expectResponseContains: "",
		},
		{
			name:                   "AddFeature",
			choice:                 AddFeature,
			detailsFlag:            "",
			expectedErr:            nil,
			expectResponseContains: "Implement the following feature(s):",
		},
		{
			name:                   "FixBug",
			choice:                 FixBug,
			detailsFlag:            "",
			expectedErr:            nil,
			expectResponseContains: "Fix the following bug(s):",
		},
		{
			name:                   "CodeCleanup",
			choice:                 CodeCleanup,
			detailsFlag:            "",
			expectedErr:            nil,
			expectResponseContains: "Cleanup and refactor this codebase.",
		},
		{
			name:                   "WriteTests",
			choice:                 WriteTests,
			detailsFlag:            "",
			expectedErr:            nil,
			expectResponseContains: "Write unit tests for this codebase.",
		},
		{
			name:                   "AddFeatureWithDetailsFlag",
			choice:                 AddFeature,
			detailsFlag:            "This is a test using detailsFlag",
			expectedErr:            nil,
			expectResponseContains: "This is a test using detailsFlag",
		},
		{
			name:        "InvalidSelection",
			choice:      99,
			detailsFlag: "",
			expectedErr: errors.New("invalid selection"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response, err := buildPrompt(fileContents, test.choice, test.detailsFlag)
			if err != nil && test.expectedErr == nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if err == nil && test.expectedErr != nil {
				t.Errorf("Expected error: %v, got nil", test.expectedErr)
			}
			if err != nil && test.expectedErr != nil && err.Error() != test.expectedErr.Error() {
				t.Errorf("Expected error: %v, got %v", test.expectedErr, err)
			}
			if test.expectResponseContains != "" && !strings.Contains(response, test.expectResponseContains) {
				t.Errorf("Expected response to contain %v, got %v", test.expectResponseContains, response)
			}
			if test.detailsFlag != "" && !strings.Contains(response, test.detailsFlag) {
				t.Errorf("Expected response to contain detailsFlag %v, got %v", test.detailsFlag, response)
			}
		})
	}
}
