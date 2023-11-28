package parsers

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/koniferous22/patched-git/git-wrapper/config"
)

func extractTargetDirectoryFromClonedRepositoryRef(cloneRepositoryRef string) (string, error) {
	// should work on ssh, https, and bare repositories
	hasDotGitSuffix := strings.HasSuffix(cloneRepositoryRef, ".git")
	shouldKeepDotGitSuffix := false
	// keep .git suffix only in case of bare repos
	if hasDotGitSuffix {
		_, err := os.Stat(cloneRepositoryRef)
		if err == nil {
			shouldKeepDotGitSuffix = true
		} else if !os.IsNotExist(err) {
			return "", fmt.Errorf("error extracting cloned directory from repository ref\n%w", err)
		}
	}
	pathFragments := strings.Split(cloneRepositoryRef, "/")
	defaultGitCloneOutput := pathFragments[len(pathFragments)-1]
	if shouldKeepDotGitSuffix {
		return defaultGitCloneOutput, nil
	}
	return strings.TrimSuffix(defaultGitCloneOutput, ".git"), nil
}

func ParseGitRepositoryPathFromArgs(appConfig config.Config, dashCFlag string, operation string, operationArgs []string) (string, error) {
	handleError := func(err error) error {
		return fmt.Errorf("error parsing git repository path from args\n%w", err)
	}
	pwd := dashCFlag
	if pwd == "" {
		pwd = "."
	}
	switch operation {
	case "init":
		// Join -C flag with optional directory arg
		gitInitDirectoryFromArgs, err := ParseGitInitDirectoryFromArgs(appConfig, operationArgs)
		if err != nil {
			return "", handleError(err)
		}
		return path.Join(pwd, gitInitDirectoryFromArgs.FoundArg), nil
	case "clone":
		// Join -C flag with inferred clone destination potentially overridden by extra clone flag
		// example - git clone ssh://some-repo clone-dest
		cloningRef, err := ParseGitCloningRefFromArgs(appConfig, operationArgs)
		if err != nil {
			return "", handleError(err)
		}
		if cloningRef.FoundArg == "" {
			return "", handleError(fmt.Errorf("cloning ref was not found in \"git clone\" args"))
		}
		expectedClonedRelativePath, err := extractTargetDirectoryFromClonedRepositoryRef(cloningRef.FoundArg)
		if err != nil {
			return "", handleError(err)
		}
		if len(cloningRef.Rest) > 0 {
			expectedClonedRelativePath = cloningRef.Rest[0]
		}
		return path.Join(pwd, expectedClonedRelativePath), nil
	default:
		return pwd, nil
	}
}
