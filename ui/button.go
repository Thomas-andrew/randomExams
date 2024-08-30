package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

type ColorButton struct {
	widget.BaseWidget

	rect    *canvas.Rectangle
	text    *widget.Label
	pressed bool
	tapFunc func()
}

func NewColorButton(title string, tp func()) *ColorButton {
	rect := canvas.NewRectangle(color.RGBA{91, 87, 117, 255})
	text := widget.NewLabel(title)
	cb := &ColorButton{
		rect:    rect,
		text:    text,
		pressed: false,
		tapFunc: tp,
	}

	cb.ExtendBaseWidget(cb)

	return cb
}

func (cb *ColorButton) CreateRenderer() fyne.WidgetRenderer {
	return &ColorButtonRender{
		widget:  cb,
		objects: []fyne.CanvasObject{cb.rect, cb.text},
	}
}

type ColorButtonRender struct {
	widget  *ColorButton
	objects []fyne.CanvasObject
}

func (c *ColorButtonRender) MinSize() fyne.Size {
	return c.widget.text.MinSize()
}

func (c *ColorButtonRender) Layout(fyne.Size) {
	c.widget.rect.Resize(c.widget.text.MinSize())
}

func (c *ColorButtonRender) Objects() []fyne.CanvasObject {
	return c.objects
}

func (c *ColorButtonRender) Refresh() {
	canvas.Refresh(c.widget)
}

func (c *ColorButtonRender) Destroy() {}

func (cb *ColorButton) Tapped(_ *fyne.PointEvent) {
	cb.pressed = !cb.pressed
	if cb.pressed {
		cb.rect.FillColor = color.RGBA{31, 0, 204, 255}
	} else {
		cb.rect.FillColor = color.RGBA{91, 87, 117, 255}
	}
	cb.tapFunc()
	canvas.Refresh(cb.rect)
}
