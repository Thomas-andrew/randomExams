package main

import (
	"fmt"
	"log"
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
			num, err := strconv.Atoi(entry.Text)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}

			// grid := container.New(layout.NewGridLayout(2))
			grid := NewDynamicGrid()
			w.SetContent(grid.container)

			for i := 0; i < num; i++ {

				path := fmt.Sprintf("./imgs/imageTest-%v.png", i)
				err := screenshoot(path)
				if err != nil {
					dialog.ShowError(err, w)
					return
				}

				image := canvas.NewImageFromFile(path)
				image.SetMinSize(fyne.NewSize(200, 100))
				image.FillMode = canvas.ImageFillContain

				button := widget.NewButton(fmt.Sprintf("retake %v", i+1), func() {
					log.Printf("tapped button %v\n", i)
					err := screenshoot(path)
					if err != nil {
						dialog.ShowError(err, w)
						return
					}
					image.Refresh()
				})
				centerButton := container.New(layout.NewCenterLayout(), button)
				// grid.Add(centerButton)

				grid.AddRow(centerButton, image)
				w.SetContent(grid.container)
			}
		},
	}
	return form
}

type IngestRow struct {
	buttons []fyne.CanvasObject
	images  []fyne.CanvasObject

	basePath string
}

func NewIngestRow(w fyne.Window) *IngestRow {
	ingest := &IngestRow{
		basePath: "./imgs/img_test",
	}
	button := widget.NewButton("add image", func() {
		ingest.AddImage(w)
	})

	ingest.buttons = append(ingest.buttons, button)

	return ingest
}

func (i *IngestRow) AddImage(w fyne.Window) {
	num := len(i.buttons)
	path := i.basePath + strconv.Itoa(num) + ".png"
	screenshoot(path)
	img := canvas.NewImageFromFile(path)
	i.images = append(i.images, img)

	button := widget.NewButton(fmt.Sprintf("retake %v", num+1), func() {
		err := screenshoot(path)
		if err != nil {
			dialog.ShowError(err, w)
		}
		img.Refresh()
	})

	i.buttons = append(i.buttons, button)
}
