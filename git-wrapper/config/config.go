package config

import (
	"fmt"
	"os"

	"github.com/caarlos0/env"
	"gopkg.in/yaml.v3"
)

type ExecutionConfig struct {
	PromptBeforeExecution bool   `yaml:"prompt,omitempty"`
	AppendRepositoryArg   bool   `yaml:"shouldAppendRepositoryArg,omitempty"`
	TargetCommand         string `yaml:"targetCommand"`
}

type OperationYamlConfig struct {
	Enabled             bool            `yaml:"enabled,omitempty"`
	SynopsisJSON        string          `yaml:"synopsisJSON"`
	Execution           ExecutionConfig `yaml:"execution"`
	FlagBlacklist       []string        `yaml:"flagBlacklist,omitempty"`
	GlobalFlagBlacklist []string        `yaml:"globalFlagBlacklist,omitempty"`
}

type GlobalYamlConfig struct {
	Enabled                     bool     `yaml:"enabled,omitempty"`
	SynopsisJSON                string   `yaml:"synopsisJSON"`
	AlwaysPromptBeforeExecution bool     `yaml:"promptBeforeExecution,omitempty"`
	FlagBlacklist               []string `yaml:"flagBlacklist,omitempty"`
}

type EnvVariables struct {
	Enabled                     bool   `env:"PATCHED_GIT_IS_ENABLED"`
	AlwaysPromptBeforeExecution bool   `env:"PATCHED_GIT_ALWAYS_PROMPT_BEFORE_EXECUTION"`
	ConfigFile                  string `env:"PATCHED_GIT_CONFIG,required"`
}

type YamlConfig struct {
	GlobalConfig   GlobalYamlConfig               `yaml:"global"`
	OperationCofig map[string]OperationYamlConfig `yaml:"operations"`
}

type Config struct {
	Env EnvVariables
	App YamlConfig
}

func loadYamlConfig(yamlConfigPath string) (*YamlConfig, error) {
	handleError := func(err error) error {
		return fmt.Errorf("error loading config.yml\n%w", err)
	}
	yamlFile, err := os.ReadFile(yamlConfigPath)
	if err != nil {
		return nil, handleError(err)
	}

	var yamlConfig YamlConfig
	err = yaml.Unmarshal(yamlFile, &yamlConfig)
	if err != nil {
		return nil, handleError(err)
	}
	return &yamlConfig, nil
}

func LoadConfig() (*Config, error) {
	handleError := func(err error) error {
		return fmt.Errorf("error loading config\n%w", err)
	}
	envVariables := EnvVariables{}
	if err := env.Parse(&envVariables); err != nil {
		return nil, handleError(fmt.Errorf("error processing env variables\n%w", err))
	}
	yamlConfig, err := loadYamlConfig(envVariables.ConfigFile)
	if err != nil {
		return nil, handleError(err)
	}
	return &Config{Env: envVariables, App: *yamlConfig}, nil
}

func (config Config) ResolveShouldRun() bool {
	return config.Env.Enabled && config.App.GlobalConfig.Enabled
}

func (config Config) ResolveShouldRunOperation(operation string) bool {
	if !config.ResolveShouldRun() {
		return false
	}
	operationConfig, ok := config.App.OperationCofig[operation]
	if !ok {
		return false
	}
	return operationConfig.Enabled
}

func (config Config) ResolveShouldPromptBeforeExecution(operation string) bool {
	if config.Env.AlwaysPromptBeforeExecution || config.App.GlobalConfig.AlwaysPromptBeforeExecution {
		return true
	}
	operationConfig, ok := config.App.OperationCofig[operation]
	if !ok {
		return false
	}
	return operationConfig.Execution.PromptBeforeExecution
}
