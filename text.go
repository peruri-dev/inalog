package inalog

import (
	"log/slog"
	"os"
)

func createTextHandler(cfg Cfg) slog.Handler {
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: cfg.Source,
		Level:     slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey {
				level := a.Value.Any().(slog.Level)
				levelLabel, exists := LevelNames[level]
				if !exists {
					levelLabel = level.String()
				}

				a.Value = slog.StringValue(levelLabel)
			}
			if a.Key == slog.MessageKey && cfg.MessageKey {
				a.Key = "message"
			}

			return a
		},
	})

	return handler
}
