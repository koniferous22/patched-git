package parsers

import (
	"testing"

	"github.com/koniferous22/patched-git/git-wrapper/config"
)

func TestParseGitRepositoryPath(t *testing.T) {
	appConfig, err := config.LoadConfig()
	if err != nil {
		t.Errorf("Error loading config")
		return
	}
	testCases := []struct {
		inputDashCFlag string
		inputOperation string
		inputArgs      []string
		expectedOutput string
		expectedErr    bool
	}{
		// git init tests
		{
			inputDashCFlag: "",
			inputOperation: "init",
			inputArgs:      []string{},
			expectedOutput: ".",
			expectedErr:    false,
		},
		{
			inputDashCFlag: "test",
			inputOperation: "init",
			inputArgs:      []string{},
			expectedOutput: "test",
			expectedErr:    false,
		},
		{
			inputDashCFlag: "",
			inputOperation: "init",
			inputArgs:      []string{"init-dir"},
			expectedOutput: "init-dir",
			expectedErr:    false,
		},
		{
			inputDashCFlag: "test",
			inputOperation: "init",
			inputArgs:      []string{"init-dir"},
			expectedOutput: "test/init-dir",
			expectedErr:    false,
		},
		// git clone tests - error
		{
			inputDashCFlag: "",
			inputOperation: "clone",
			inputArgs:      []string{},
			expectedOutput: "",
			expectedErr:    true,
		},
		// git clone tests -
		{
			inputDashCFlag: "",
			inputOperation: "clone",
			inputArgs:      []string{"ref"},
			expectedOutput: "ref",
			expectedErr:    false,
		},
		{
			inputDashCFlag: "test",
			inputOperation: "clone",
			inputArgs:      []string{"ref"},
			expectedOutput: "test/ref",
			expectedErr:    false,
		},
		{
			inputDashCFlag: "",
			inputOperation: "clone",
			inputArgs:      []string{"ref", "override"},
			expectedOutput: "override",
			expectedErr:    false,
		},
		{
			inputDashCFlag: "test",
			inputOperation: "clone",
			inputArgs:      []string{"ref", "override"},
			expectedOutput: "test/override",
			expectedErr:    false,
		},
	}

	for _, tc := range testCases {
		output, err := ParseGitRepositoryPathFromArgs(*appConfig, tc.inputDashCFlag, tc.inputOperation, tc.inputArgs)

		if (err != nil) != tc.expectedErr {
			formattedVerb := "not to receive"
			if tc.expectedErr {
				formattedVerb = "to receive"
			}
			t.Errorf("For input %v, expected %s error, but got: %v", tc.inputArgs, formattedVerb, err)
		}

		if tc.expectedOutput != output {
			t.Errorf("For input %v, expected output: %v, but got: %v", tc.inputArgs, tc.expectedOutput, output)
		}
	}
}
