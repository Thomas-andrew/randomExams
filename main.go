package main

import (
	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.New()
	w := a.NewWindow("Exercise Checklist")

	w.SetContent(makeGUI())
	w.ShowAndRun()
}

// func screenshoot(path string) error {
// 	args := []string{"-s", "-m", "10", path}
// 	cmd := exec.Command("maim", args...)
//
// 	err := cmd.Run()
// 	if err != nil {
// 		return fmt.Errorf("[screenshoot] maim finished with error: %w", err)
// 	}
//
// 	return nil
// }
