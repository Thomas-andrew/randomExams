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

type IngestRow struct {
	buttons []fyne.CanvasObject
	images  []fyne.CanvasObject

	buttonsCont *fyne.Container
	imagesCont  *fyne.Container

	basePath string
	numList  int
}

func NewIngestRow(pg string, w fyne.Window) *IngestRow {
	btnCont := container.New(layout.NewVBoxLayout())
	imgCont := container.New(layout.NewVBoxLayout())

	ingest := &IngestRow{
		basePath:    "./imgs/img_test-" + pg,
		buttonsCont: btnCont,
		imagesCont:  imgCont,
	}
	button := widget.NewButton("add image", func() {
		ingest.AddImage(w)
	})

	ingest.buttons = append(ingest.buttons, button)
	btnCont.Add(button)

	return ingest
}

func (i *IngestRow) AddImage(w fyne.Window) {
	num := len(i.buttons)
	path := i.basePath + strconv.Itoa(num) + ".png"
	screenshoot(path)
	img := canvas.NewImageFromFile(path)
	img.SetMinSize(fyne.NewSize(200, 100))
	img.FillMode = canvas.ImageFillContain
	i.images = append(i.images, img)
	i.imagesCont.Add(img)

	button := widget.NewButton(fmt.Sprintf("retake %v-%v", i.numList, num), func() {
		err := screenshoot(path)
		if err != nil {
			dialog.ShowError(err, w)
		}
		img.Refresh()
	})

	i.buttons = append(i.buttons, button)
	i.buttonsCont.Add(button)
}
