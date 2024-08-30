package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
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

func (g *GUI) StartPage() {
	buttonExams := widget.NewButton(
		"prova aleatoria",
		func() {
			// go to exams
		},
	)

	buttonIngest := widget.NewButton(
		"adicionar exercicios",
		func() {
			makeIngestForm(g)
		},
	)

	cont := container.NewCenter(
		container.NewVBox(
			buttonExams,
			buttonIngest,
		),
	)

	g.Window.SetContent(cont)
}
