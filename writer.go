package inalog

import (
	"context"
	"log/slog"
	"os"
	"runtime"
	"time"
)

type writer struct {
	logImpl    *slog.Logger
	ctx        context.Context
	skipCaller int
}

func newWriter(w writer) writer {
	return w
}

func shouldPrint(currentLv slog.Level) bool {
	if currentLv == LevelNotice {
		return true
	}

	if currentLv > slog.LevelWarn {
		return true
	}

	getEnv := os.Getenv("INALOG_LOG_LEVEL")
	getLevel := slog.LevelDebug
	switch getEnv {
	case "INFO", "info":
		{
			getLevel = slog.LevelInfo
		}
	case "WARN", "warn":
		{
			getLevel = slog.LevelWarn
		}
	}

	return currentLv >= getLevel
}

func (w writer) logHandling(lvl slog.Level, msg string, attr ...any) {
	if !shouldPrint(lvl) {
		return
	}

	var pcs [1]uintptr
	// add skip caller here, due to custom writer add 3 caller layer
	w.skipCaller += 3
	attrs := []slog.Attr{}

	if w.ctx == nil {
		w.ctx = context.Background()
	} else {
		for _, hook := range hooks {
			attrs = append(attrs, hook(w.ctx)...)
		}
	}

	runtime.Callers(w.skipCaller, pcs[:])
	r := slog.NewRecord(time.Now(), lvl, msg, pcs[0])
	r.Add(attr...)
	r.AddAttrs(attrs...)
	w.logImpl.Handler().Handle(
		w.ctx,
		r,
	)
}

func (h InalogHandler) Handle(ctx context.Context, r slog.Record) error {
	attrs := []slog.Attr{}

	for _, hook := range hooks {
		attrs = append(attrs, hook(ctx)...)
	}

	r.AddAttrs(attrs...)

	return h.Handler.Handle(ctx, r)
}
