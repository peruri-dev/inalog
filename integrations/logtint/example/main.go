package main

import (
	"context"
	"log/slog"
	"os"
	"strconv"

	"github.com/peruri-dev/inalog"
	"github.com/peruri-dev/inalog/integrations/logtint"
)

type ctxAppName struct{}

func appNameExtract(ctx context.Context) []slog.Attr {
	attrs := []slog.Attr{}
	appName := ctx.Value(ctxAppName{})
	if appName != nil {
		attrs = append(attrs, slog.String("application-name", appName.(string)))
	}

	return attrs
}

func main() {
	isJsonLog, _ := strconv.ParseBool(os.Getenv("JSON_LOG"))

	cfg := inalog.Cfg{
		Source:     true,
		TextLog:    !isJsonLog,
		CustomFunc: logtint.CreateTintHandler(),
		MessageKey: true,
	}

	inalog.Init(cfg)

	inalog.AddHook(appNameExtract)

	slog.Info("Information", slog.String("key", "value"))
	slog.Debug("Debug", slog.String("key", "value"))
	slog.Warn("Warning", slog.String("key", "value"))

	ctx := context.WithValue(context.Background(), ctxAppName{}, "example-inalog")
	slog.ErrorContext(ctx, "Error", slog.String("key", "value"))

	inalog.Log().Notice("POST /pingz")
	inalog.LogWith(inalog.WithCfg{Ctx: ctx}).Notice("GET /healthz")
}
