package layouts

import (
	"fyne.io/fyne/v2"
)

type ExamLayout struct{}

func (t *ExamLayout) Layout(objs []fyne.CanvasObject, size fyne.Size) {
	if len(objs) != 4 {
		panic("wrong number of objects passed to TopBarExamLayout")
	}
	winW, winH := size.Components()
	timerH := objs[0].Size().Height
	barH := objs[1].Size().Height
	buttonH := objs[2].Size().Height

	h := max(timerH, barH, buttonH, 50)

	// timer
	objs[0].Resize(fyne.NewSize(winW/3, h))
	objs[0].Move(fyne.NewPos(winW/6, 5))

	// bar
	objs[1].Resize(fyne.NewSize(winW/3, h))
	objs[1].Move(fyne.NewPos(winW/3, 0))

	// button
	objs[2].Resize(fyne.NewSize(winW/3, h))
	objs[2].Move(fyne.NewPos(2*winW/3, 0))

	// scroll
	objs[3].Resize(fyne.NewSize(winW, winH-h))
	objs[3].Move(fyne.NewPos(0, h))
}

func (t *ExamLayout) MinSize(objs []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(100, 100)
}

func max(heights ...float32) float32 {
	var maxH float32 = 0
	for _, height := range heights {
		if height > maxH {
			maxH = height
		}
	}
	return maxH
}
