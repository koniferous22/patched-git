package parsers

import (
	"testing"

	"github.com/koniferous22/patched-git/git-wrapper/config"
)

func comparePositionalArgsResults(a, b PositionalArgsParsingResult) bool {
	if a.FoundArg != b.FoundArg {
		return false
	}

	if len(a.Rest) != len(b.Rest) {
		return false
	}

	for i := range a.Rest {
		if a.Rest[i] != b.Rest[i] {
			return false
		}
	}

	for i := range a.EscapedFlags {
		if a.EscapedFlags[i] != b.EscapedFlags[i] {
			return false
		}
	}
	return true
}

func TestParseOperationFromArgs(t *testing.T) {
	appConfig, err := config.LoadConfig()
	if err != nil {
		t.Errorf("Error loading config")
		return
	}
	testCases := []struct {
		inputArgs      []string
		expectedOutput *PositionalArgsParsingResult
		expectedErr    error
	}{
		{
			inputArgs: []string{"-C", "clone", "init"},
			expectedOutput: &PositionalArgsParsingResult{
				FoundArg: "init",
				EscapedFlags: map[string]string{
					"-C": "clone",
				},
				Rest: []string{},
			},
			expectedErr: nil,
		},
		{
			inputArgs: []string{"-Ca", "clone", "init"},
			expectedOutput: &PositionalArgsParsingResult{
				FoundArg:     "clone",
				EscapedFlags: map[string]string{},
				Rest:         []string{"init"},
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		output, err := ParseGitOperationFromArgs(*appConfig, tc.inputArgs)

		if err != tc.expectedErr {
			t.Errorf("For input %v, expected error: %v, but got: %v", tc.inputArgs, tc.expectedErr, err)
		}

		if !comparePositionalArgsResults(*tc.expectedOutput, *output) {
			t.Errorf("For input %v, expected output: %v, but got: %v", tc.inputArgs, tc.expectedOutput, output)
		}
	}
}
