package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/koniferous22/patched-git/git-wrapper/config"
	"github.com/koniferous22/patched-git/git-wrapper/parsers"
	"github.com/koniferous22/patched-git/git-wrapper/utils"
)

func main() {
	handleError := func(err error) {
		fmt.Fprintf(os.Stderr, "%s%s%s\n", utils.ColorRed, err.Error(), utils.Reset)
		os.Exit(1)
	}
	config, err := config.LoadConfig()
	if err != nil {
		handleError(err)
	}
	if !config.ResolveShouldRun() {
		return
	}
	args := os.Args[1:]
	if utils.ContainsAnyOf(args, config.App.GlobalConfig.FlagBlacklist) {
		return
	}
	parseOperationOutput, err := parsers.ParseGitOperationFromArgs(*config, args)
	if err != nil {
		handleError(err)
	}
	parsedGitCommand := parseOperationOutput.FoundArg
	operationArgs := parseOperationOutput.Rest
	if !config.ResolveShouldRunOperation(parsedGitCommand) {
		return
	}

	if utils.ContainsAnyOf(args, config.App.OperationCofig[parsedGitCommand].GlobalFlagBlacklist) ||
		utils.ContainsAnyOf(operationArgs, config.App.OperationCofig[parsedGitCommand].FlagBlacklist) {
		return
	}
	targetCommand := config.App.OperationCofig[parsedGitCommand].Execution.TargetCommand
	shouldExecuteTarget := true
	if config.ResolveShouldPromptBeforeExecution(parsedGitCommand) {
		shouldExecuteTarget, err = runYesNo(fmt.Sprintf("Do you want to run \"%s\" for \"%s\" git command", targetCommand, parsedGitCommand))
		if err != nil {
			handleError(err)
		}
	}
	shouldAppendRepositoryArg := config.App.OperationCofig[parsedGitCommand].Execution.AppendRepositoryArg
	if shouldExecuteTarget {
		command := targetCommand
		var commandArgs []string
		if shouldAppendRepositoryArg {
			dashCGlobalArg := ""
			if dashCValue, ok := parseOperationOutput.EscapedFlags["-C"]; ok {
				dashCGlobalArg = dashCValue
			}
			repositoryArg, err := parsers.ParseGitRepositoryPathFromArgs(*config, dashCGlobalArg, parsedGitCommand, operationArgs)
			if err != nil {
				handleError(err)
			}
			repositoryArgAbsPath, err := filepath.Abs(repositoryArg)
			if err != nil {
				handleError(err)
			}
			commandArgs = append(commandArgs, repositoryArgAbsPath)
		}
		cmd := exec.Command(command, commandArgs...)
		stdout, err := cmd.Output()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		fmt.Println(string(stdout))
	}
}
