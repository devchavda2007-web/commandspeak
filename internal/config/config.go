package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	DefaultDeployPlatform string            `json:"defaultDeployPlatform"`
	ProjectPath           string            `json:"projectPath"`
	Commands              map[string]string `json:"commands"`
}

const configFileName = ".commandspeakrc.json"

// GetConfigPath returns the path to the config file in the user's home directory.
func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, configFileName), nil
}

// LoadConfig reads the config file and returns a Config struct.
func LoadConfig() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Return default config if file does not exist
			return createDefaultConfig(), nil
		}
		return nil, err
	}

	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	// Ensure the commands map is initialized if missing
	if cfg.Commands == nil {
		cfg.Commands = createDefaultConfig().Commands
	}

	return &cfg, nil
}

// SaveConfig saves the given Config struct to the config file.
func SaveConfig(cfg *Config) error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

func createDefaultConfig() *Config {
	return &Config{
		DefaultDeployPlatform: "vercel",
		ProjectPath:           ".",
		Commands: map[string]string{
			"DEPLOY":         "git add . && git commit -m \"update\" && git push && {{platform}} --prod",
			"PUSH":           "git add . && git commit -m \"{{message}}\" && git push",
			"CLEAN_BRANCHES": "git branch --merged | grep -v \"\\*\" | grep -v \"master\" | grep -v \"main\" | xargs -r -n 1 git branch -d",
			"RUN_TESTS":      "npm test && npm run build",
		},
	}
}
