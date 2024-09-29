package conf

import (
	"encoding/json"
	"fmt"
	"os"
)

type Configuration struct {
	DBHost       string `json:"DBHost"`
	DBPort       int    `json:"DBPort"`
	DBUser       string `json:"DBUser"`
	DBPassword   string `json:"DBPassword"`
	DBName       string `json:"DBName"`
	ServiceName  string `json:"ServiceName"`
	ServicePort  int    `json:"ServicePort"`
	ObserverPort int    `json:"ObserverPort"`
}

func LoadConfiguration(configPath string) *Configuration {
	config := Configuration{}
	if configPath == "" {
		fmt.Println("Configuration file not found")
	}
	if fileExists(configPath) {
		if file, err := os.ReadFile(configPath); err == nil {
			errJson := json.Unmarshal(file, &config)
			if errJson != nil {
				fmt.Println("Failed to read configuration file, not able to unmarshal", errJson)
			}
		} else {
			fmt.Println("Failed to load configurations, Not able to open config file", err)
		}
	}
	return &config
}

func fileExists(filePath string) bool {
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
