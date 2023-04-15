package wfc

import (
	"github.com/gotk3/gotk3/cairo"
)

var (
	tiles map[string]*cairo.Surface
)

func (w *wfc) loadTiles(path string) error {
	config, err := LoadConfig(path)
	if err != nil {
		return err
	}

	surface, err := cairo.NewSurfaceFromPNG(config.Path)
	if err != nil {
		return err
	}

	tiles = make(map[string]*cairo.Surface)
	for _, tile := range config.Tiles {
		tiles[tile.Key] = surface.CreateForRectangle(tile.Left, tile.Top, tile.Width, tile.Height)
	}

	return nil
}
