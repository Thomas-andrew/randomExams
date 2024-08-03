package main

import (
	"fyne.io/fyne/v2"
)

type rExLayout struct {
	top, left, content fyne.CanvasObject
}

const sideWidth = 100

func newrExLayout(top, left, content fyne.CanvasObject) fyne.Layout {
	return &rExLayout{top: top, left: left, content: content}
}

func (l *rExLayout) Layout(_ []fyne.CanvasObject, size fyne.Size) {
	// top geometry
	topHeight := l.top.MinSize().Height
	l.top.Resize(fyne.NewSize(size.Width, topHeight))

	// left geometry
	l.left.Move(fyne.NewPos(0, topHeight))
	l.left.Resize(fyne.NewSize(sideWidth, size.Height-topHeight))

	l.content.Move(fyne.NewPos(sideWidth, topHeight))
	l.content.Resize(fyne.NewSize(size.Width-sideWidth, size.Height-topHeight))
}

func (l *rExLayout) MinSize(_ []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(100, 100)
}
