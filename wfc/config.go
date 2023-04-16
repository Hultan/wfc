package wfc

import (
	"encoding/json"
	"os"
)

type Config struct {
	Path  string `json:"path"`
	Tiles []struct {
		Key    string  `json:"key"`
		Top    float64 `json:"top"`
		Left   float64 `json:"left"`
		Width  float64 `json:"width"`
		Height float64 `json:"height"`
	} `json:"tiles"`
}

// LoadConfig : Loads the configuration file
func LoadConfig(path string) (*Config, error) {
	// Open config file
	configFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(configFile *os.File) {
		_ = configFile.Close()
	}(configFile)

	// Parse the JSON document
	var config Config
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
