package main

import "fyne.io/fyne/v2"

type ingestRowLayout struct{}

func NewIngestRowLayout() *ingestRowLayout {
	return &ingestRowLayout{}
}

const padding float32 = 20

func (l *ingestRowLayout) Layout(objs []fyne.CanvasObject, size fyne.Size) {
	if len(objs) != 2 {
		return
	}

	leftWidth := size.Width * 0.1
	rightWidth := size.Width * 0.9

	// buttons
	objs[0].Resize(fyne.NewSize(leftWidth, size.Height))
	objs[0].Move(fyne.NewPos(0, size.Height/2))

	// imgs
	objs[1].Resize(fyne.NewSize(rightWidth-padding, size.Height-padding))
	objs[1].Move(fyne.NewPos(leftWidth+padding, padding))
}

func (l *ingestRowLayout) MinSize(objs []fyne.CanvasObject) fyne.Size {
	if len(objs) != 2 {
		return fyne.NewSize(0, 0)
	}

	leftMin := objs[0].MinSize()
	rightMin := objs[1].MinSize()

	return fyne.NewSize(leftMin.Width+rightMin.Width, fyne.Max(leftMin.Height, rightMin.Height))
}
