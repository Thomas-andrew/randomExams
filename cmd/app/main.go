package main

import (
	"log/slog"
	"os"
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
	LogLevel.Set(slog.LevelDebug)
	Logger.Info("------------------------------------- application start --------------------------------------")

	gui := makeGUI()
	gui.window.ShowAndRun()
}
