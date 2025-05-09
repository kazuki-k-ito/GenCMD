package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	APIKey string `yaml:"api_key"`
	Model  string `yaml:"model"`
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".gen_cmd_config.yaml"), nil
}

func saveConfig(apiKey string, model string) error {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	config := Config{
		APIKey: apiKey,
		Model:  model,
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
