package data

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

type GUI struct {
	Window fyne.Window
	App    fyne.App
}

func MakeGUI() *GUI {
	app := app.New()
	win := app.NewWindow("random exercises")

	return &GUI{
		App:    app,
		Window: win,
	}
}

func (g *GUI) ShowAndRun() {
	g.Window.ShowAndRun()
}
