package wfc

import (
	"encoding/json"
	"os"
)

type Config struct {
	Path  string `json:path`
	Tiles []struct {
		Key    string  `json:key`
		Top    float64 `json:top`
		Left   float64 `json:left`
		Width  float64 `json:width`
		Height float64 `json:height`
	} `json:tiles`
}

// LoadConfig : Loads the configuration file
func LoadConfig(path string) (*Config, error) {
	// Make sure the file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

	// Open config file
	configFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	// Parse the JSON document
	var config *Config
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	if err != nil {
		return nil, err
	}

	_ = configFile.Close()

	return config, nil
}

// Save : Saves the configuration file
func (tileMap *Config) Save() error {
	// Open config file
	configFile, err := os.OpenFile(configPath, os.O_TRUNC|os.O_WRONLY, 0644)

	// Handle errors
	if err != nil {
		return err
	}

	// Create JSON from config object
	data, err := json.MarshalIndent(tileMap, "", "\t")
	if err != nil {
		return err
	}

	// Write the data
	_, err = configFile.Write(data)
	if err != nil {
		return err
	}

	_ = configFile.Close()

	return nil
}
