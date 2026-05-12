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

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	fmt.Println("homeDir :", homeDir)

	return homeDir + "/" + configFileName, nil
}

func Read() (Config, error) {
	filePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, nil
	}

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return Config{}, err
	}
	fmt.Println("File Content: ", string(fileContent))

	cfg := Config{}
	err = json.Unmarshal(fileContent, &cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}


func SetUser(cfg *Config)