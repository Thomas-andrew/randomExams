package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type ColorButton struct {
	widget.BaseWidget

	overlay *canvas.Rectangle
	text    *widget.Label
	pressed bool
}

func NewColorButton(title string) *ColorButton {
	rect := canvas.NewRectangle(color.RGBA{91, 87, 117, 255})
	text := widget.NewLabel(title)
	cb := &ColorButton{
		overlay: rect,
		text:    text,
		pressed: false,
	}

	cb.ExtendBaseWidget(cb)

	return cb
}

func (cb *ColorButton) CreateRenderer() fyne.WidgetRenderer {
	content := container.New(layout.NewPaddedLayout(), cb.overlay, cb.text)
	return widget.NewSimpleRenderer(content)
}

func (cb *ColorButton) Tapped(_ *fyne.PointEvent) {
	cb.pressed = !cb.pressed
	if cb.pressed {
		cb.overlay.FillColor = color.RGBA{31, 0, 204, 255}
	} else {
		cb.overlay.FillColor = color.RGBA{91, 87, 117, 255}
	}
	canvas.Refresh(cb.overlay)
}
