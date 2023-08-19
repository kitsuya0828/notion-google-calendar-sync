package main

import (
	"log"
	"os"

	"github.com/Kitsuya0828/notion-googlecalendar-sync/run"
	"golang.org/x/exp/slog"
)

func main() {
	opt := &slog.HandlerOptions{
		// AddSource: true,
		Level: slog.LevelDebug,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opt))

	err := run.Run(logger)
	if err != nil {
		log.Fatal(err)
	}
}
