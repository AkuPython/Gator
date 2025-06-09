package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (config *Config) SetUser(user string) error {
	configPath, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("Could not get home dir: %v", err)
	}

	config.CurrentUserName = user

	return writeJSON(config, configPath)

}

func writeJSON(j *Config, f string) error {
	file, err := os.Create(f)
	if err != nil {
		return fmt.Errorf("Could not Write Config: %v", err)
	}
	defer file.Close()
	
	// Use json.Encoder to write JSON efficiently
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(j); err != nil {
		return fmt.Errorf("Encoder / writing issue: %v", err)
	}
	return nil
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("Could not get home dir: %v", err)
	}
	return home + "/" + configFileName, nil
}

func Read() (Config, error) {
	configPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, fmt.Errorf("Could not get home dir: %v", err)
	}
	file, err := os.Open(configPath)
	if err != nil {
		return Config{}, fmt.Errorf("Could not open config file: %v", err)
	}
	defer file.Close()
	
	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return Config{}, fmt.Errorf("Error unmarshaling JSON: %v", err)
	}

	return config, nil

}
