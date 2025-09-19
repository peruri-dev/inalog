package logtint

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/peruri-dev/inalog"
)

func CreateTintHandler() func(cfg inalog.Cfg) slog.Handler {

	return func(cfg inalog.Cfg) slog.Handler {
		w := os.Stderr
		handler := tint.NewHandler(w, &tint.Options{
			AddSource:  cfg.Source,
			Level:      slog.LevelDebug,
			TimeFormat: time.Kitchen,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.LevelKey {
					if level, ok := a.Value.Any().(slog.Level); ok {
						if levelLabel, exists := inalog.ShortLevelNames[level]; exists {
							return tint.Attr(13, slog.String(a.Key, levelLabel))
						}
					}
					//a.Value = slog.StringValue(levelLabel)
				}

				return a
			},
		})
		return handler
	}
	
}
