package wfc

import (
	"math/rand"
	"time"

	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gtk"
)

const (
	configPath = "/home/per/code/wfc/roads.json"
	TileSize   = 64
	LeftMargin = 10
	TopMargin  = 10
	Height     = 200
	Width      = 200
)

type gridCell struct {
	tileKey     string
	isCollapsed bool
}

type wfc struct {
	da     *gtk.DrawingArea
	random *rand.Rand
	grid   [Height][Width]gridCell
	tiles  map[string]*cairo.Surface
}

func newWFC(da *gtk.DrawingArea) *wfc {
	return &wfc{
		da:     da,
		random: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (w *wfc) setup() error {
	w.da.Connect("draw", w.onDraw)

	err := w.loadTiles(configPath)
	if err != nil {
		return err
	}

	w.regenerate()

	return nil
}

func (w *wfc) loadTiles(path string) error {
	config, err := LoadConfig(path)
	if err != nil {
		return err
	}

	surface, err := cairo.NewSurfaceFromPNG(config.Path)
	if err != nil {
		return err
	}

	w.tiles = make(map[string]*cairo.Surface)
	for _, tile := range config.Tiles {
		w.tiles[tile.Key] = surface.CreateForRectangle(tile.Left, tile.Top, tile.Width, tile.Height)
	}

	return nil
}

func (w *wfc) onDraw(da *gtk.DrawingArea, ctx *cairo.Context) {
	ctx.SetSourceRGBA(0.8, 0.8, 0.8, 0.5)
	ctx.Rectangle(0, 0, float64(da.GetAllocatedWidth()), float64(da.GetAllocatedHeight()))
	ctx.Fill()

	for y := 0; y < Height; y++ {
		for x := 0; x < Width; x++ {
			xx, yy := float64(x)*TileSize+LeftMargin, float64(y)*TileSize+TopMargin
			w.drawTile(ctx, w.tiles[w.grid[y][x].tileKey], xx, yy)
		}
	}
}

func (w *wfc) drawTile(ctx *cairo.Context, surface *cairo.Surface, x, y float64) {
	ctx.SetSourceSurface(surface, x, y)
	ctx.Paint()
}

func (w *wfc) generateWorld() {

	for y := 0; y < Height; y++ {
		for x := 0; x < Width; x++ {
			connections := "????"

			if y > 0 && w.grid[y-1][x].isCollapsed {
				connections = replaceCharInString(connections, string(w.grid[y-1][x].tileKey[2]), 0)
			} else if y == 0 {
				connections = replaceCharInString(connections, "0", 0)
			}
			if x < Width-1 && w.grid[y][x+1].isCollapsed {
				connections = replaceCharInString(connections, string(w.grid[y][x+1].tileKey[3]), 1)
			} else if x == Width-1 {
				connections = replaceCharInString(connections, "0", 1)
			}
			if y < Height-1 && w.grid[y+1][x].isCollapsed {
				connections = replaceCharInString(connections, string(w.grid[y+1][x].tileKey[0]), 2)
			} else if y == Height-1 {
				connections = replaceCharInString(connections, "0", 2)
			}
			if x > 0 && w.grid[y][x-1].isCollapsed {
				connections = replaceCharInString(connections, string(w.grid[y][x-1].tileKey[1]), 3)
			} else if x == 0 {
				connections = replaceCharInString(connections, "0", 3)
			}

			// Pick a random matching tile
			key := w.pickRandomMatchingTile(connections)

			w.grid[y][x].tileKey = key
			w.grid[y][x].isCollapsed = true
		}
	}
}

func (w *wfc) pickRandomMatchingTile(connections string) string {
	keys := getKeys(w.tiles)

	// Remove non-matching tiles
	for i := len(keys) - 1; i >= 0; i-- {
		if !w.isKeyValid(connections, keys[i]) {
			keys[i] = keys[len(keys)-1]
			keys = keys[:len(keys)-1]
		}
	}

	return keys[w.random.Intn(len(keys))]
}

func (w *wfc) isKeyValid(connections string, key string) bool {
	for i := 0; i < len(connections); i++ {
		if (connections[i] == '0' || connections[i] == '1') && connections[i] != key[i] {
			return false
		}
	}

	return true
}

func (w *wfc) regenerate() {
	// Clear the grid
	w.grid = [Height][Width]gridCell{}

	w.generateWorld()
	w.da.QueueDraw()
}
