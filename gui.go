package main

import (
	"fmt"
	"os/exec"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
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
			var imgWidgets []fyne.CanvasObject
			for i := 0; i < num; i++ {
				path := fmt.Sprintf("./imgs/imageTest-%v.png", i)
				err := screenshoot(path)
				if err != nil {
					dialog.ShowError(err, w)
					return
				}

				image := canvas.NewImageFromFile(path)
				image.SetMinSize(fyne.NewSize(200, 100))
				imgWidgets = append(imgWidgets, image)
			}

			grid := container.New(layout.NewGridLayout(1), imgWidgets...)
			w.SetContent(grid)
		},
	}
	return form
}

func screenshoot(path string) error {
	args := []string{"-s", "-m", "10", path}
	cmd := exec.Command("maim", args...)

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("[screenshoot] maim finished with error: %w", err)
	}

	return nil
}
