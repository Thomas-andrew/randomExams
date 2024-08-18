package main

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func makeExIngest(w fyne.Window) fyne.CanvasObject {
	entry := widget.NewEntry()
	entry.SetPlaceHolder("Enter number of screenshoots")

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "number", Widget: entry},
		},
		OnSubmit: func() {
			numEntries, err := strconv.Atoi(entry.Text)
			if err != nil {
				dialog.ShowError(err, w)
			}

			vertList := container.New(layout.NewVBoxLayout())

			scrollVertList := container.NewScroll(vertList)

			for i := 1; i <= numEntries; i++ {
				ingestRowData := NewIngestRow(i, w)
				ingestRowData.AddImage(w)

				ingestRow := container.New(
					NewIngestRowLayout(),
					ingestRowData.buttons,
					ingestRowData.images,
				)
				vertList.Add(ingestRow)
				w.SetContent(scrollVertList)
			}
		},
	}
	return form
	// fyne.Layout
}

type ingestData struct {
	images  *fyne.Container
	buttons *fyne.Container

	path string
	num  int
	id   int
}

func NewIngestRow(id int, w fyne.Window) *ingestData {
	images := container.NewVBox()
	buttons := container.NewVBox()

	ingest := &ingestData{
		images:  images,
		buttons: buttons,

		path: "./imgs/img_test-pog",
		num:  0,
		id:   id,
	}

	addButton := widget.NewButton(
		"Add image",
		func() {
			ingest.AddImage(w)
		},
	)

	ingest.buttons.Add(addButton)

	return ingest
}

func (g *ingestData) AddImage(w fyne.Window) {
	g.num += 1
	path := g.path + "-" + strconv.Itoa(g.id) + strconv.Itoa(g.num) + ".png"
	err := screenshoot(path)
	if err != nil {
		dialog.ShowError(err, w)
	}

	img := canvas.NewImageFromFile(path)
	img.SetMinSize(fyne.NewSize(700, 500))
	img.FillMode = canvas.ImageFillContain
	g.images.Add(img)

	retakeButton := widget.NewButton(
		fmt.Sprintf("retake %v", g.num),
		func() {
			err := screenshoot(path)
			if err != nil {
				dialog.ShowError(err, w)
			}

			img.Refresh()
		},
	)

	g.buttons.Add(retakeButton)
}
