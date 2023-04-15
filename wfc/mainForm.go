package wfc

import (
	"fmt"
	"os"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/softteam/framework"
)

const applicationTitle = "wave function collapse"
const applicationVersion = "version 0.01"
const applicationCopyRight = "©SoftTeam AB, 2022"

type MainForm struct {
	Window  *gtk.ApplicationWindow
	builder *framework.GtkBuilder
}

// NewMainForm : Creates a new MainForm object
func NewMainForm() *MainForm {
	mainForm := new(MainForm)
	return mainForm
}

var w *wfc

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
	m.Window = m.builder.GetObject("main_window").(*gtk.ApplicationWindow)

	// Set up main window
	m.Window.SetApplication(app)
	title := fmt.Sprintf("%s - %s", applicationTitle, applicationVersion)
	m.Window.SetTitle(title)
	m.Window.Maximize()

	// Hook up events
	m.Window.Connect("destroy", m.Window.Close)
	m.Window.Connect("key-press-event", m.onKeyDown)

	// Quit button
	button := m.builder.GetObject("main_window_quit_button").(*gtk.ToolButton)
	button.Connect("clicked", m.Window.Close)

	// Status bar
	statusBar := m.builder.GetObject("main_window_status_bar").(*gtk.Statusbar)
	statusBar.Push(statusBar.GetContextId("wfc"), title)

	// Drawing area
	da := m.builder.GetObject("drawing_area").(*gtk.DrawingArea)
	w = newWFC(da)
	err = w.setup()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	// Menu
	m.setupMenu()

	// Show the main window
	m.Window.ShowAll()
}

func (m *MainForm) setupMenu() {
	menuQuit := m.builder.GetObject("menu_file_quit").(*gtk.MenuItem)
	menuQuit.Connect("activate", m.Window.Close)
}
func (m *MainForm) onKeyDown(_ *gtk.ApplicationWindow, e *gdk.Event) {
	ke := gdk.EventKeyNewFromEvent(e)

	if ke.KeyVal() == gdk.KEY_F5 {
		w.regenerate()
	}
}
