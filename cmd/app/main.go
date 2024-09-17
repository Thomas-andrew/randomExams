package main

import (
	"log/slog"
	"os"

	"github.com/Twintat/randomExams/data"
	"github.com/Twintat/randomExams/ui"
)

var (
	Logger   *slog.Logger
	LogLevel = new(slog.LevelVar)
)

func main() {
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	Logger = slog.New(
		slog.NewTextHandler(
			// os.Stdout,
			logFile,
			&slog.HandlerOptions{
				Level: LogLevel,
			},
		),
	)
	slog.SetDefault(Logger)
	LogLevel.Set(slog.LevelDebug)
	Logger.Info("------------------------------------- application start --------------------------------------")

	gui := data.MakeGUI()
	ui.StartPage(gui)
	gui.ShowAndRun()
}
