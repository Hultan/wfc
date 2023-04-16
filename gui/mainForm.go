package gui

import (
	"fmt"
	"log"
	"os"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/softteam/framework"
	"github.com/hultan/wfc/wfc"
)

const applicationTitle = "wave function collapse"
const applicationVersion = "version 0.01"
const applicationCopyRight = "©SoftTeam AB, 2022"

type MainForm struct {
	window  *gtk.ApplicationWindow
	builder *framework.GtkBuilder
	wfc     *wfc.Wfc
}

// NewMainForm : Creates a new MainForm object
func NewMainForm() *MainForm {
	mainForm := new(MainForm)
	return mainForm
}

// OpenMainForm : Opens the MainForm window
func (m *MainForm) OpenMainForm(app *gtk.Application) {
	// Initialize gtk
	gtk.Init(&os.Args)

	// Create a new softBuilder
	fw := framework.NewFramework()
	builder, err := fw.Gtk.CreateBuilder("main.glade")
	if err != nil {
		panic(err)
	}
	m.builder = builder

	// Get the main window from the glade file
	m.window = m.builder.GetObject("main_window").(*gtk.ApplicationWindow)

	// Set up main window
	m.window.SetApplication(app)
	title := fmt.Sprintf("%s - %s - %s", applicationTitle, applicationVersion, applicationCopyRight)
	m.window.SetTitle(title)
	m.window.Maximize()

	// Hook up events
	m.window.Connect("destroy", m.window.Close)
	m.window.Connect("key-press-event", m.onKeyDown)

	// Quit button
	button := m.builder.GetObject("main_window_quit_button").(*gtk.ToolButton)
	button.Connect("clicked", m.window.Close)

	// Status bar
	statusBar := m.builder.GetObject("main_window_status_bar").(*gtk.Statusbar)
	statusBar.Push(statusBar.GetContextId("wfc"), title)

	// Drawing area
	da := m.builder.GetObject("drawing_area").(*gtk.DrawingArea)
	m.wfc, err = wfc.NewWFC(da)
	if err != nil {
		log.Fatal(err)
	}

	// Menu
	m.setupMenu()

	// Show the main window
	m.window.ShowAll()
}

func (m *MainForm) setupMenu() {
	menuQuit := m.builder.GetObject("menu_file_quit").(*gtk.MenuItem)
	menuQuit.Connect("activate", m.window.Close)
}
func (m *MainForm) onKeyDown(_ *gtk.ApplicationWindow, e *gdk.Event) {
	ke := gdk.EventKeyNewFromEvent(e)

	switch ke.KeyVal() {
	case gdk.KEY_F5:
		m.wfc.Generate()
	case gdk.KEY_q, gdk.KEY_Q:
		m.window.Close()
	}
}
