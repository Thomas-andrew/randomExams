package ui

import (
	"fyne.io/fyne/v2"
)

type newSetLayout struct{}

func (n *newSetLayout) Layout(objs []fyne.CanvasObject, size fyne.Size) {
	if len(objs) != 2 {
		return
	}
	winW, winH := size.Components()
	// set the button size
	var butH float32 = 50

	// move button
	objs[1].Move(fyne.NewPos(0, winH-butH))
	objs[1].Resize(fyne.NewSize(winW, butH))
	// resize tree
	objs[0].Resize(fyne.NewSize(winW, winH-butH))
	objs[0].Move(fyne.NewPos(0, 0))
}

func (n *newSetLayout) MinSize(objs []fyne.CanvasObject) fyne.Size {
	if len(objs) != 2 {
		return fyne.Size{}
	}
	return fyne.NewSize(
		objs[0].MinSize().Width+objs[1].MinSize().Width,
		objs[0].MinSize().Height+objs[1].MinSize().Height,
	)
}
