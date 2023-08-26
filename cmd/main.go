package main

import (
	"log"
	"os"

	"github.com/Kitsuya0828/notion-google-calendar-sync/run"
	"golang.org/x/exp/slog"
)

func main() {
	opt := &slog.HandlerOptions{
		// AddSource: true,
		Level: slog.LevelDebug,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opt))
	slog.SetDefault(logger)

	err := run.Run()
	if err != nil {
		log.Fatal(err)
	}
}
