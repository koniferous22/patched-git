package parsers

import (
	"fmt"

	"github.com/koniferous22/patched-git/git-wrapper/config"
	"github.com/koniferous22/patched-git/git-wrapper/utils"
)

type PositionalArgsParsingResult struct {
	FoundArg     string
	EscapedFlags map[string]string
	Rest         []string
}

// Find 1st positional argument that isn't preceeded by one of found flags
func parsePositionalArgWithEscapedFlagArguments(args []string, flags []string) PositionalArgsParsingResult {
	var rest []string
	result := PositionalArgsParsingResult{
		FoundArg:     "",
		EscapedFlags: make(map[string]string),
	}
	for i, arg := range args {
		if !utils.IsFlag(arg) {
			if i > 0 && utils.Contains(flags, args[i-1]) {
				result.EscapedFlags[args[i-1]] = arg
				continue
			}
			if i+1 < len(args) {
				rest = append(rest, args[i+1:]...)
			}
			result.FoundArg = arg
			result.Rest = rest
			return result
		}
	}
	return result
}

func parsePositionalArgsWithGeneratedSynopsis(synopsisPath string, args []string, errMessage string) (*PositionalArgsParsingResult, error) {
	handleError := func(err error) error {
		return fmt.Errorf("%s\n%w", errMessage, err)
	}
	// Example: Load parser output for `man git` synopsis
	synopsisFormat, err := config.ParseSynopsisJson(synopsisPath)
	if err != nil {
		return nil, handleError(err)
	}
	// Example - Load space delimited flags, that affect positional arg resolution
	// example - "git -C clone init" would incorrectly resolve to "clone" if not escaped
	spaceDelimitedFlagOpts, err := synopsisFormat.GetSpaceDelimitedFlagOptions()
	if err != nil {
		return nil, handleError(err)
	}
	positionalArgParsingResult := parsePositionalArgWithEscapedFlagArguments(args, *spaceDelimitedFlagOpts)
	return &positionalArgParsingResult, nil
}

func ParseGitOperationFromArgs(appConfig config.Config, args []string) (*PositionalArgsParsingResult, error) {
	return parsePositionalArgsWithGeneratedSynopsis(appConfig.App.GlobalConfig.SynopsisJSON, args, "error resolving git command")
}

func ParseGitInitDirectoryFromArgs(appConfig config.Config, gitInitArgs []string) (*PositionalArgsParsingResult, error) {
	return parsePositionalArgsWithGeneratedSynopsis(appConfig.App.OperationCofig["init"].SynopsisJSON, gitInitArgs, "error resolving git init directory")
}

func ParseGitCloningRefFromArgs(appConfig config.Config, gitInitArgs []string) (*PositionalArgsParsingResult, error) {
	return parsePositionalArgsWithGeneratedSynopsis(appConfig.App.OperationCofig["clone"].SynopsisJSON, gitInitArgs, "error resolving git cloning ref")
}
