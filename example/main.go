package main

import (
	"log/slog"

	"github.com/peruri-dev/inalog"
)

func main() {
	inalog.Init(inalog.Cfg{
		Source: true,
		Tinted: true,
	})

	slog.Info("Information", slog.String("key", "value"))
	slog.Debug("Debug", slog.String("key", "value"))
	slog.Warn("Warning", slog.String("key", "value"))
	slog.Error("Error", slog.String("key", "value"))
}
