package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func makeGUI() fyne.CanvasObject {
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.HomeIcon(), func() {}),
	)
	left := widget.NewLabel("Left")
	content := makeContent()
	objs := []fyne.CanvasObject{toolbar, left, content}
	return container.New(newrExLayout(toolbar, left, content), objs...)
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

		button := NewColorButton(item.Label)

		itemContainer := container.New(
			layout.NewHBoxLayout(),
			button,
			img,
		)
		content.Add(itemContainer)
	}
	return content
}
