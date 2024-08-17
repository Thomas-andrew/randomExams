package main

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
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

			grid := NewDynamicGrid()
			w.SetContent(grid.container)

			for i := 0; i < num; i++ {

				ingest := NewIngestRow("Ex"+strconv.Itoa(i+1), w)
				ingest.AddImage(w)

				grid.AddRow(ingest.buttonsCont, ingest.imagesCont)
				w.SetContent(grid.container)
			}
		},
	}
	return form
}
