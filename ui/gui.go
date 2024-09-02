package ui

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Twintat/randomExams/data"
)

func StartPage(g *data.GUI) {
	buttonExams := widget.NewButton(
		"prova aleatoria",
		func() {
			startRandomExam(g)
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
