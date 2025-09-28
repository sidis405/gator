package config

import (
	"encoding/json"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbUrl           string `json:"db_url,omitempty"`
	CurrentUserName string `json:"current_user_name,omitempty"`
}

func (c *Config) SetUser(userName string) error {
	c.CurrentUserName = userName

	return write(c)
}

func write(c *Config) error {
	jsonContents, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(getConfigFilePath(), jsonContents, 0777)
}

func Read() (Config, error) {
	configContents, err := os.ReadFile(getConfigFilePath())
	if err != nil {
		return Config{}, err
	}

	var config Config
	if err := json.Unmarshal(configContents, &config); err != nil {
		return Config{}, err
	}

	return config, nil
}

func getConfigFilePath() string {
	homePath, _ := os.UserHomeDir()
	return homePath + "/" + configFileName
}
