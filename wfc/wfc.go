package wfc

import (
	"math/rand"
	"time"

	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gtk"
)

const (
	mapHeight = 10
	mapWidth  = 10
)

type gridCell struct {
	tileKey     string
	isCollapsed bool
}

type Wfc struct {
	da                    *gtk.DrawingArea
	random                *rand.Rand
	grid                  [mapHeight][mapWidth]gridCell
	tiles                 map[string]*cairo.Surface
	tileHeight, tileWidth float64
}

func NewWFC(da *gtk.DrawingArea, configPath string) (*Wfc, error) {
	w := &Wfc{
		da:     da,
		random: rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	// Hook up the draw signal
	w.da.Connect("draw", w.onDraw)

	// Load the tiles
	err := w.loadTiles(configPath)
	if err != nil {
		return nil, err
	}

	// Generate the "world"
	w.Generate()

	return w, nil
}

// Generate generates a new "world"
func (w *Wfc) Generate() {
	// Clear the grid
	w.grid = [mapHeight][mapWidth]gridCell{}

	// Generate a new "world" and draw it
	w.generateWorld()
	w.da.QueueDraw()
}

//
// Private functions
//

func (w *Wfc) loadTiles(path string) error {
	// Load the config file
	config, err := LoadConfig(path)
	if err != nil {
		return err
	}

	err = config.Validate()
	if err != nil {
		return err
	}

	w.tileWidth = config.TileWidth
	w.tileHeight = config.TileHeight

	// Load the tile map
	surface, err := cairo.NewSurfaceFromPNG(config.TileMap.Path)
	if err != nil {
		return err
	}

	// Add tiles to the map
	w.tiles = make(map[string]*cairo.Surface)
	for _, tile := range config.TileMap.Tiles {
		w.tiles[tile.Key] = surface.CreateForRectangle(tile.Left, tile.Top, tile.Width, tile.Height)
	}

	return nil
}

func (w *Wfc) onDraw(da *gtk.DrawingArea, ctx *cairo.Context) {
	// Draw the gray background
	ctx.SetSourceRGBA(0.8, 0.8, 0.8, 0.5)
	ctx.Rectangle(0, 0, float64(da.GetAllocatedWidth()), float64(da.GetAllocatedHeight()))
	ctx.Fill()

	// Draw the "world"
	for y := 0; y < mapHeight; y++ {
		for x := 0; x < mapWidth; x++ {
			xx, yy := float64(x)*w.tileWidth, float64(y)*w.tileHeight
			ctx.SetSourceSurface(w.tiles[w.grid[y][x].tileKey], xx, yy)
			ctx.Paint()
		}
	}
}

func (w *Wfc) generateWorld() {

	for y := 0; y < mapHeight; y++ {
		for x := 0; x < mapWidth; x++ {
			pattern := "????"

			// Up
			if y == 0 {
				pattern = replaceCharInString(pattern, "0", 0)
			} else if w.grid[y-1][x].isCollapsed {
				pattern = replaceCharInString(pattern, string(w.grid[y-1][x].tileKey[2]), 0)
			}
			// Right
			if x == mapWidth-1 {
				pattern = replaceCharInString(pattern, "0", 1)
			} else if w.grid[y][x+1].isCollapsed {
				pattern = replaceCharInString(pattern, string(w.grid[y][x+1].tileKey[3]), 1)
			}
			// Down
			if y == mapHeight-1 {
				pattern = replaceCharInString(pattern, "0", 2)
			} else if w.grid[y+1][x].isCollapsed {
				pattern = replaceCharInString(pattern, string(w.grid[y+1][x].tileKey[0]), 2)
			}
			// Left
			if x == 0 {
				pattern = replaceCharInString(pattern, "0", 3)
			} else if w.grid[y][x-1].isCollapsed {
				pattern = replaceCharInString(pattern, string(w.grid[y][x-1].tileKey[1]), 3)
			}

			w.grid[y][x].tileKey = w.pickRandomMatchingTile(pattern)
			w.grid[y][x].isCollapsed = true
		}
	}
}

func (w *Wfc) pickRandomMatchingTile(pattern string) string {
	keys := getKeys(w.tiles)

	// Remove non-matching tiles
	for i := len(keys) - 1; i >= 0; i-- {
		if !w.isKeyValid(pattern, keys[i]) {
			keys[i] = keys[len(keys)-1]
			keys = keys[:len(keys)-1]
		}
	}

	return keys[w.random.Intn(len(keys))]
}

func (w *Wfc) isKeyValid(pattern string, key string) bool {
	// Does the key match the pattern?
	for i := 0; i < len(pattern); i++ {
		if (pattern[i] == '0' || pattern[i] == '1') && pattern[i] != key[i] {
			return false
		}
	}

	return true
}
