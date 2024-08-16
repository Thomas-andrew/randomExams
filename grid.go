package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type GridLayout struct {
	supTextWidth float32
}

func NewCustomGridLayout() *GridLayout {
	return &GridLayout{
		supTextWidth: 0,
	}
}

var padding float32 = 10

func (g *GridLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	log.Printf("-------------- ingest page size: {H: %v, W: %v} -------------------\n", size.Height, size.Width)
	lineHeight := size.Height / float32(len(objects)/2)
	xPos := float32(0)
	yPos := float32(0)

	g.supTextWidth = size.Width * 0.08

	log.Printf(
		"data: {lineHeight: %v, supTextWidth: %v, xPos: %v, yPos: %v}\n",
		lineHeight, g.supTextWidth, xPos, yPos,
	)

	for i, obj := range objects {
		if i%2 == 0 {
			// first column
			obj.Resize(fyne.NewSize(g.supTextWidth, lineHeight))
			obj.Move(fyne.NewPos(xPos, yPos))
			log.Printf("button\t{%v,%v} -> size{h: %v, w: %v}\n", i/2+1, i%2, obj.Size().Height, obj.Size().Width)
		} else {
			// second column
			obj.Resize(fyne.NewSize(size.Width-g.supTextWidth-padding, lineHeight-padding))
			obj.Move(fyne.NewPos(xPos+g.supTextWidth, yPos))
			log.Printf("image\t{%v,%v} -> size{h: %v, w: %v}\n", i/2+1, i%2, obj.Size().Height, obj.Size().Width)
			yPos += lineHeight
		}
	}
}

func (g *GridLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	var supTextWidth float32 = 0
	for i, obj := range objects {
		if i%2 == 0 {
			width := obj.Size().Width
			if width >= supTextWidth {
				supTextWidth = width
			}
		}
	}

	minSize := fyne.NewSize(supTextWidth, 0)

	for i, obj := range objects {
		if i%2 == 1 {
			minSize.Height += obj.MinSize().Height
		} else {
			minSize.Width = fyne.Max(minSize.Width, obj.MinSize().Width)
		}
	}
	return minSize
}

type DynamicGrid struct {
	container *fyne.Container
	layout    *GridLayout
}

func NewDynamicGrid() *DynamicGrid {
	lyt := NewCustomGridLayout()
	cont := container.New(lyt)
	return &DynamicGrid{
		container: cont,
		layout:    lyt,
	}
}

func (d *DynamicGrid) AddRow(first fyne.CanvasObject, second fyne.CanvasObject) {
	d.container.Add(first)
	d.container.Add(second)
}

func (d *DynamicGrid) objects() []fyne.CanvasObject {
	return d.container.Objects
}

func (d *DynamicGrid) Refresh() {
	d.container.Refresh()
}

func (d *DynamicGrid) Hide() {}
