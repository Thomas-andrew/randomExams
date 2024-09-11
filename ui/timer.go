package ui

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2/widget"
	"github.com/Twintat/randomExams/data"
)

type Timer struct {
	duration string
	bar      *widget.ProgressBar
	timer    *widget.Label
	button   *widget.Button
	state    TimerState
}

type TimerState int

const (
	PreStart TimerState = iota
	Running
	Finished
)

func newTimer(gui *data.GUI, duration string) *Timer {
	timer := &Timer{
		duration: duration,
		bar:      widget.NewProgressBar(),
		timer:    widget.NewLabel("00:00:00"),
		state:    PreStart,
	}
	ctx, cancel := context.WithCancel(context.Background())
	var button *widget.Button
	button = widget.NewButton(
		"começar prova",
		func() {
			switch timer.state {
			case PreStart:
				timer.state = Running
				button.SetText("terminar prova")
				go updateTimer(ctx, timer)
				return
			case Running:
				timer.state = Finished
				button.SetText("menu principal")
				cancel()
				return
			case Finished:
				StartPage(gui)
				return
			}
		},
	)
	timer.button = button
	return timer
}

func getDuration(str string) (time.Duration, int) {
	better := strings.ReplaceAll(str, " ", "")

	hoursTxt, rest, exist := strings.Cut(better, "h")
	if !exist {
		rest = hoursTxt
		hoursTxt = ""
	}
	minutesTxt, _, _ := strings.Cut(rest, "m")

	if minutesTxt == "" {
		minutesTxt = "0"
	}
	if hoursTxt == "" {
		hoursTxt = "0"
	}

	hours, err := strconv.Atoi(hoursTxt)
	if err != nil {
		panic("[getDuration] error converting hoursTxt to hours")
	}
	minutes, err := strconv.Atoi(minutesTxt)
	if err != nil {
		panic("[getDuration] error converting minutesTxt to minutes")
	}

	duration := time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute
	mins := hours*60 + minutes
	return duration, mins
}

func updateTimer(ctx context.Context, t *Timer) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	duration, totalMin := getDuration(t.duration)
	timer := time.NewTimer(duration)
	defer timer.Stop()

	ratio := 1 / float64(totalMin*60)
	var progress float64 = 0
	for {
		select {
		case <-ticker.C:
			// update timer
			progress += 1
			hours, rest := divmod(int64(progress), 3600)
			minutes, seconds := divmod(rest, 60)
			t.bar.SetValue(progress * ratio)
			t.timer.SetText(
				fmt.Sprintf(
					"%02d:%02d:%02d",
					hours, minutes, seconds,
				),
			)
		case <-timer.C:
			// end exam
			t.timer.SetText("-- prova terminada --")
			t.bar.SetValue(1)
			return
		case <-ctx.Done():
			// end exam
			t.timer.SetText("-- prova terminada --")
			t.bar.SetValue(1)
			return
		}
	}
}

func divmod(up, down int64) (quo, rem int64) {
	quo = up / down
	rem = up % down
	return
}
