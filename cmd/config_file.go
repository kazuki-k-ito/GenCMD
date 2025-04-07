package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	APIKey          string `yaml:"api_key"`
	Shell           string `yaml:"shell"`
	Model           string `yaml:"model"`
	OperationSystem string `yaml:"operation_system"`
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".gen_cmd_config.yaml"), nil
}

func saveConfig(apiKey string, shell string, model string, operationSystem string) error {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	config := Config{
		APIKey:          apiKey,
		Shell:           shell,
		Model:           model,
		OperationSystem: operationSystem,
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	err = os.WriteFile(configFilePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func loadConfig() (Config, error) {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return Config{}, nil
	}

	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func init() {
	var err error
	_, err = loadConfig()
	if err != nil {
		fmt.Println("Failed to load config:", err)
	}
}
