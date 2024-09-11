package ui

import (
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

type ExerciseExamLayout struct{}

func (e *ExerciseExamLayout) Layout(objs []fyne.CanvasObject, size fyne.Size) {
	// images
	winW, winH := size.Components()
	slog.Debug(
		"size given for exercise exam layout",
		"h", winH,
		"w", winW,
	)

	var totalH float32 = 0
	for _, obj := range objs {
		img, ok := obj.(*canvas.Image)
		if !ok {
			return
		}
		aspectRatio := img.Aspect()
		sz := fyne.NewSize(winW, winW/aspectRatio)
		img.Resize(sz)
		img.SetMinSize(sz)
		img.Move(fyne.NewPos(0, totalH))
		totalH += sz.Height
		slog.Debug(
			"images to layout",
			"name", img.File,
			"h", img.Size().Height,
			"w", img.Size().Width,
			"ratio", aspectRatio,
		)
	}
}

func (e *ExerciseExamLayout) MinSize(objs []fyne.CanvasObject) fyne.Size {
	w, h := float32(0), float32(0)
	for _, obj := range objs {
		childSize := obj.MinSize()
		w += childSize.Width
		h += childSize.Height
	}
	return fyne.NewSize(w, h)
}
