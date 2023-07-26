package main

import (
	"encoding/json"
	"fmt"
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gtk"
	"os"
)

type JSONMain struct {
	MapWidth   int         `json:"mapWidth"`
	MapHeight  int         `json:"mapHeight"`
	TileWidth  float64     `json:"tileWidth"`
	TileHeight float64     `json:"tileHeight"`
	TileMap    JSONTileMap `json:"tileMap"`
}

type JSONTileMap struct {
	Path  string     `json:"path"`
	Tiles []JSONTile `json:"tiles"`
}
type JSONTile struct {
	Key    string `json:"key"`
	Top    int    `json:"top"`
	Left   int    `json:"left"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

const (
	tileMapPath = "/home/per/code/wfc/assets/world.png"
	tileWidth   = 16
	tileHeight  = 16
)

var (
	tileMap              JSONMain
	pathEntry            *gtk.Entry
	keyEntry             *gtk.Entry
	pic                  *gtk.DrawingArea
	surface, currentTile *cairo.Surface
	currentX, currentY   float64
	current              int
)

func main() {
	gtk.Init(nil)

	w, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	panicOnError(err)
	w.SetTitle("Json generator")
	w.SetDefaultSize(400, 300)

	box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	panicOnError(err)
	w.Add(box)

	pathEntry, err = gtk.EntryNew()
	panicOnError(err)
	pathEntry.SetText(tileMapPath)
	box.Add(pathEntry)

	btn, err := gtk.ButtonNewWithLabel("Load")
	panicOnError(err)
	btn.Connect("clicked", load)
	box.Add(btn)

	pic, err = gtk.DrawingAreaNew()
	panicOnError(err)
	pic.SetHExpand(true)
	pic.SetVExpand(true)
	pic.Connect("draw", draw)
	box.Add(pic)

	keyEntry, err = gtk.EntryNew()
	panicOnError(err)
	keyEntry.SetText("")
	box.Add(keyEntry)

	btn, err = gtk.ButtonNewWithLabel("Next")
	panicOnError(err)
	btn.Connect("clicked", next)
	box.Add(btn)

	btn, err = gtk.ButtonNewWithLabel("Done")
	panicOnError(err)
	btn.Connect("clicked", done)
	box.Add(btn)

	w.Connect("destroy", closeWindow)
	w.ShowAll()

	gtk.Main()
}

func load() {
	path, err := pathEntry.GetText()
	panicOnError(err)

	surface, err = cairo.NewSurfaceFromPNG(path)
	panicOnError(err)

	currentTile = surface.CreateForRectangle(0, 0, tileWidth, tileHeight)
	current = 0
	pic.QueueDraw()
	keyEntry.GrabFocus()
}

func next() {
	key, err := keyEntry.GetText()
	panicOnError(err)

	if key != "" {
		t := JSONTile{
			Key:    key,
			Left:   int(currentX),
			Top:    int(currentY),
			Width:  tileWidth,
			Height: tileHeight,
		}
		tileMap.TileMap.Tiles = append(tileMap.TileMap.Tiles, t)
		keyEntry.SetText("")
	}

	width, height := float64(surface.GetWidth()), float64(surface.GetHeight())

	currentX += tileWidth
	if currentX >= width {
		currentX = 0
		currentY += tileHeight
		fmt.Println()
	}

	if currentY >= height {
		// We are done
		done()
		return
	}

	fmt.Printf("Rect(%v,%v)\n", currentX, currentY)
	currentTile = surface.CreateForRectangle(currentX, currentY, tileWidth, tileHeight)
	current += 1
	keyEntry.GrabFocus()
	pic.QueueDraw()

	return
}

func done() {
	tileMap.MapWidth = 20
	tileMap.MapHeight = 20
	tileMap.TileWidth = tileWidth
	tileMap.TileHeight = tileHeight
	tileMap.TileMap.Path = tileMapPath

	f, err := os.Create("map.json")
	panicOnError(err)
	json.NewEncoder(f).Encode(tileMap)
}

func draw(da *gtk.DrawingArea, ctx *cairo.Context) {
	ctx.SetSourceSurface(surface, 0, 0)
	ctx.Paint()

	ctx.SetSourceRGBA(50, 50, 50, 255)
	ctx.Rectangle(currentX, currentY, tileWidth, tileHeight)
	ctx.Stroke()
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func closeWindow() {
	gtk.MainQuit()
}
