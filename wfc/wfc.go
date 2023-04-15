package wfc

import (
	"math/rand"
	"time"

	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gtk"
)

const (
	TilesPath  = "/home/per/code/wfc/assets/tiles.png"
	TileSize   = 64
	LeftMargin = 0
	TopMargin  = 0
	Height     = 200
	Width      = 200
)

type gridCell struct {
	tileKey     string
	isCollapsed bool
}

type position struct {
	x, y int
}

type wfc struct {
	da *gtk.DrawingArea
}

var grid [Height][Width]gridCell
var random *rand.Rand

func newWFC(da *gtk.DrawingArea) *wfc {
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
	return &wfc{da}
}

func (w *wfc) setup() error {
	w.da.Connect("draw", w.onDraw)

	surface, err := w.loadSurface(TilesPath)
	if err != nil {
		return err
	}
	w.createTiles(surface)

	w.regenerate()

	return nil
}

func (w *wfc) onDraw(da *gtk.DrawingArea, ctx *cairo.Context) {
	ctx.SetSourceRGBA(0.8, 0.8, 0.8, 0.5)
	ctx.Rectangle(0, 0, float64(da.GetAllocatedWidth()), float64(da.GetAllocatedHeight()))
	ctx.Fill()

	for y := 0; y < Height; y++ {
		for x := 0; x < Width; x++ {
			xx, yy := float64(x)*TileSize+LeftMargin, float64(y)*TileSize+TopMargin
			cell := grid[y][x]
			w.drawTile(ctx, tiles[cell.tileKey], xx, yy)
		}
	}
}

func (w *wfc) drawTile(ctx *cairo.Context, surface *cairo.Surface, x, y float64) {
	ctx.SetSourceSurface(surface, x, y)
	ctx.Paint()
}

func (w *wfc) pickRandomEmptySquare() (int, int) {
	var toPick []position

	for y := 0; y < Height; y++ {
		for x := 0; x < Width; x++ {
			if !grid[y][x].isCollapsed {
				toPick = append(toPick, position{x, y})
			}
		}
	}

	if len(toPick) == 0 {
		return -1, -1
	}

	t := random.Intn(len(toPick))
	return toPick[t].x, toPick[t].y
}

func (w *wfc) generateWorld() {

	for y := 0; y < Height; y++ {
		for x := 0; x < Width; x++ {
			connections := "????"

			if y > 0 && grid[y-1][x].isCollapsed {
				connections = replacePartOfString(connections, string(grid[y-1][x].tileKey[2]), 0)
			} else if y == 0 {
				connections = replacePartOfString(connections, "0", 0)
			}
			if x < Width-1 && grid[y][x+1].isCollapsed {
				connections = replacePartOfString(connections, string(grid[y][x+1].tileKey[3]), 1)
			} else if x == Width-1 {
				connections = replacePartOfString(connections, "0", 1)
			}
			if y < Height-1 && grid[y+1][x].isCollapsed {
				connections = replacePartOfString(connections, string(grid[y+1][x].tileKey[0]), 2)
			} else if y == Height-1 {
				connections = replacePartOfString(connections, "0", 2)
			}
			if x > 0 && grid[y][x-1].isCollapsed {
				connections = replacePartOfString(connections, string(grid[y][x-1].tileKey[1]), 3)
			} else if x == 0 {
				connections = replacePartOfString(connections, "0", 3)
			}

			// Pick a random matching tile
			key := w.pickRandomMatchingTile(connections)

			grid[y][x].tileKey = key
			grid[y][x].isCollapsed = true
		}
	}
}

func (w *wfc) pickRandomMatchingTile(connections string) string {
	keys := getKeys(tiles)

	// Remove non-matching tiles
	for i := len(keys) - 1; i >= 0; i-- {
		if !w.isKeyValid(connections, keys[i]) {
			keys[i] = keys[len(keys)-1]
			keys = keys[:len(keys)-1]
		}
	}

	if len(keys) == 0 {
		panic("No valid tiles")
	}

	t := random.Intn(len(keys))
	return keys[t]
}

func (w *wfc) isKeyValid(connections string, key string) bool {
	if len(connections) != len(key) {
		panic("len(connections) != len(key)")
	}

	for i := 0; i < len(connections); i++ {
		if (connections[i] == '0' || connections[i] == '1') && connections[i] != key[i] {
			return false
		}
	}

	return true
}

func (w *wfc) regenerate() {
	// Clear the grid
	grid = [Height][Width]gridCell{}

	w.generateWorld()
	w.da.QueueDraw()
}
