package wfc

import (
	"encoding/json"
	"errors"
	"os"
)

type Config struct {
	MapWidth   int     `json:"mapWidth"`
	MapHeight  int     `json:"mapHeight"`
	TileWidth  float64 `json:"tileWidth"`
	TileHeight float64 `json:"tileHeight"`
	TileMap    struct {
		Path  string `json:"path"`
		Tiles []struct {
			Key      string  `json:"key"`
			Top      float64 `json:"top"`
			Left     float64 `json:"left"`
			Width    float64 `json:"width"`
			Height   float64 `json:"height"`
			Priority int     `json:"priority"`
		} `json:"tiles,omitempty"`
	} `json:"tileMap,omitempty"`
	Tile []struct {
		Path string `json:"path"`
		Key  string `json:"key"`
	} `json:"tile,omitempty"`
}

var errorInvalidTileSize = errors.New("invalid tile size")

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

func (c *Config) Validate() error {
	width := c.TileMap.Tiles[0].Width
	height := c.TileMap.Tiles[0].Height

	for _, tile := range c.TileMap.Tiles {
		if tile.Width != width || tile.Height != height {
			return errorInvalidTileSize
		}
	}
	// for _, tile := range c.Tile {
	// 	if tile.Width != width || tile.Height != height {
	// 		return errorInvalidTileSize
	// 	}
	// }

	return nil
}
