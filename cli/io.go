package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

// Function to check if a file exists
func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

// Function to check if a directory exists and create it if it doesn't
func ensureDirectoryExists(dirPath string) {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, 0755) // Create the directory with 0755 permissions
		e(err)
		log.Println("Directory created:", dirPath)
	}
}

// Config represents the structure of your configuration
type Config struct {
	ServerUrl   string `yaml:"ServerUrl"`
	AccessToken string `yaml:"AccessToken"`
}

// Function to read configurations from a YAML file
func readConfigFromFile(filePath string) Config {
	var config Config
	fileContent, err := os.ReadFile(filePath)
	e(err, fmt.Sprintf("error reading config file: %v", err))
	err = yaml.Unmarshal(fileContent, &config)
	e(err, fmt.Sprintf("error unmarshaling YAML: %v", err))
	return config
}

// Function to write configurations to a YAML file
func writeConfigToFile(config Config, filePath string) {
	configBytes, err := yaml.Marshal(config)
	e(err, fmt.Sprintf("error marshaling config to YAML: %v", err))
	err = os.WriteFile(filePath, configBytes, 0644)
	e(err)
}
