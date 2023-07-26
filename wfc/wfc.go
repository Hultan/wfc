package wfc

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gtk"
)

type Surface struct {
	priority int
	surface  *cairo.Surface
}

type gridCell struct {
	tileKey     string
	isCollapsed bool
}

type Wfc struct {
	da                    *gtk.DrawingArea
	random                *rand.Rand
	grid                  [][]gridCell
	tiles                 map[string]Surface
	mapHeight, mapWidth   int
	tileHeight, tileWidth float64
}

var keySize int

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
	w.grid = w.getNewMap()

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

	w.mapWidth = config.MapWidth
	w.mapHeight = config.MapHeight
	w.tileWidth = config.TileWidth
	w.tileHeight = config.TileHeight

	// Load the tile map
	surface, err := cairo.NewSurfaceFromPNG(config.TileMap.Path)
	if err != nil {
		return err
	}

	// Add tiles to the map
	w.tiles = make(map[string]Surface)
	keySize = len(config.TileMap.Tiles[0].Key)
	for _, tile := range config.TileMap.Tiles {
		if len(tile.Key) != keySize {
			panic(fmt.Sprintf("Invalid key size : %s\n", tile.Key))
		}
		w.tiles[tile.Key] = Surface{
			priority: tile.Priority,
			surface:  surface.CreateForRectangle(tile.Left, tile.Top, tile.Width, tile.Height),
		}
	}

	return nil
}

func (w *Wfc) onDraw(da *gtk.DrawingArea, ctx *cairo.Context) {
	// Draw the gray background
	ctx.SetSourceRGBA(0.8, 0.8, 0.8, 0.5)
	ctx.Rectangle(0, 0, float64(da.GetAllocatedWidth()), float64(da.GetAllocatedHeight()))
	ctx.Fill()

	// Draw the "world"
	for y := 0; y < w.mapHeight; y++ {
		for x := 0; x < w.mapWidth; x++ {
			xx, yy := float64(x)*w.tileWidth, float64(y)*w.tileHeight
			ctx.SetSourceSurface(w.tiles[w.grid[y][x].tileKey].surface, xx, yy)
			ctx.Paint()
			//
			//ctx.SetSourceRGBA(200, 200, 200, 255)
			//ctx.Rectangle(xx, yy, w.tileWidth, w.tileHeight)
			//ctx.Stroke()
		}
	}
}

func (w *Wfc) generateWorld() {
	replace := strings.Repeat("0", keySize/4)

	for y := 0; y < w.mapHeight; y++ {
		for x := 0; x < w.mapWidth; x++ {
			pattern := strings.Repeat("?", keySize)

			// Up
			if y == 0 {
				pattern = replaceCharInString(pattern, replace, 0)
			} else if w.grid[y-1][x].isCollapsed {
				tile := w.grid[y-1][x]
				pattern = replaceCharInString(pattern, getKeyPart(tile.tileKey, 2), 0)
			}
			// Right
			if x == w.mapWidth-1 {
				pattern = replaceCharInString(pattern, replace, 1)
			} else if w.grid[y][x+1].isCollapsed {
				tile := w.grid[y][x+1]
				pattern = replaceCharInString(pattern, getKeyPart(tile.tileKey, 3), 1)
			}
			// Down
			if y == w.mapHeight-1 {
				pattern = replaceCharInString(pattern, replace, 2)
			} else if w.grid[y+1][x].isCollapsed {
				tile := w.grid[y+1][x]
				pattern = replaceCharInString(pattern, getKeyPart(tile.tileKey, 0), 2)
			}
			// Left
			if x == 0 {
				pattern = replaceCharInString(pattern, replace, 3)
			} else if w.grid[y][x-1].isCollapsed {
				tile := w.grid[y][x-1]
				pattern = replaceCharInString(pattern, getKeyPart(tile.tileKey, 1), 3)
			}

			key := w.pickRandomMatchingTile(pattern)
			if key == "INVALID KEY!" {
				return
			}
			w.grid[y][x].tileKey = key
			w.grid[y][x].isCollapsed = true
		}
		fmt.Println("-------------------------")
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

	var keyList []string
	for _, key := range keys {
		tile := w.tiles[key]
		for p := 0; p < tile.priority+1; p++ {
			keyList = append(keyList, key)
		}
	}

	if len(keyList) == 0 {
		fmt.Printf("No valid keys for pattern %s\n", pattern)
		return "INVALID KEY!"
	}

	r := keyList[w.random.Intn(len(keyList))]
	fmt.Printf("Picked square: %s\n", r)
	return r
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

func (w *Wfc) getNewMap() [][]gridCell {
	var m [][]gridCell

	for y := 0; y < w.mapHeight; y++ {
		m = append(m, []gridCell{})
		for x := 0; x < w.mapWidth; x++ {
			m[y] = append(m[y], gridCell{})
		}
	}

	return m
}
