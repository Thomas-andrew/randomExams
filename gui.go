package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type GUI struct {
	window fyne.Window
	app    fyne.App
}

func makeGUI() *GUI {
	a := app.New()
	w := a.NewWindow("Random Exercise")
	gui := &GUI{
		window: w,
		app:    a,
	}

	gui.startPage()

	return gui
}

func (g *GUI) startPage() {
	button1 := widget.NewButton("adicionar exercicios", func() {
		ingest := makeExIngest(g)
		g.window.SetContent(ingest)
	})
	button2 := widget.NewButton("fazer teste", func() {
		log.Println("pog")
	})

	addBook := widget.NewButton(
		"adicionar livro",
		func() {
			ingestBook := makeAddBook(g)
			g.window.SetContent(ingestBook)
		},
	)

	buttons := container.NewVBox(
		button1,
		button2,
		addBook,
	)

	g.window.SetContent(container.NewCenter(buttons))
}

type Checklist struct {
	Label     string
	ImagePath string
	checked   bool
}

func makeContent() fyne.CanvasObject {
	items := []Checklist{
		{"ex 1", "./imgs/elon-v1-chap.01-1.png", false},
		{"ex 2", "./imgs/elon-v1-chap.01-2.png", false},
		{"ex 3", "./imgs/elon-v1-chap.01-3.png", false},
		{"ex 4", "./imgs/elon-v1-chap.01-4.png", false},
	}

	content := container.New(layout.NewVBoxLayout())

	for _, item := range items {
		img := canvas.NewImageFromFile(item.ImagePath)
		img.FillMode = canvas.ImageFillOriginal

		button := NewColorButton(item.Label, func() {})

		itemContainer := container.New(
			layout.NewHBoxLayout(),
			button,
			img,
		)
		content.Add(itemContainer)
	}
	return content
}
