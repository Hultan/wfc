package wfc

import (
	"github.com/gotk3/gotk3/cairo"
)

var tiles map[string]*cairo.Surface

func (w *wfc) loadSurface(path string) (*cairo.Surface, error) {
	surface, err := cairo.NewSurfaceFromPNG(path)
	if err != nil {
		return nil, err
	}

	return surface, nil
}

func (w *wfc) createTiles(surface *cairo.Surface) {
	tiles = map[string]*cairo.Surface{
		"0000": surface.CreateForRectangle(250, 252, TileSize, TileSize),
		"0001": surface.CreateForRectangle(250, 92, TileSize, TileSize),
		"0010": surface.CreateForRectangle(250, 172, TileSize, TileSize),
		"0011": surface.CreateForRectangle(250, 12, TileSize, TileSize),
		"0100": surface.CreateForRectangle(90, 252, TileSize, TileSize),
		"0101": surface.CreateForRectangle(90, 92, TileSize, TileSize),
		"0110": surface.CreateForRectangle(90, 172, TileSize, TileSize),
		"0111": surface.CreateForRectangle(90, 12, TileSize, TileSize),
		"1000": surface.CreateForRectangle(170, 252, TileSize, TileSize),
		"1001": surface.CreateForRectangle(170, 92, TileSize, TileSize),
		"1010": surface.CreateForRectangle(170, 172, TileSize, TileSize),
		"1011": surface.CreateForRectangle(170, 12, TileSize, TileSize),
		"1100": surface.CreateForRectangle(10, 252, TileSize, TileSize),
		"1101": surface.CreateForRectangle(10, 92, TileSize, TileSize),
		"1110": surface.CreateForRectangle(10, 172, TileSize, TileSize),
		"1111": surface.CreateForRectangle(10, 12, TileSize, TileSize),
	}
}
