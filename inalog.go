package inalog

import (
	"log/slog"
	"os"
)

type Cfg struct {
	Source     bool
	TextLog    bool
	CustomFunc func(Cfg) slog.Handler
	MessageKey bool
}

const (
	LevelDebug  = slog.LevelDebug
	LevelInfo   = slog.LevelInfo
	LevelNotice = slog.Level(1)
	LevelError  = slog.LevelError
	LevelFatal  = slog.Level(12)
)

var LevelNames = map[slog.Leveler]string{
	LevelNotice: "NOTICE",
	LevelFatal:  "FATAL",
}
var ShortLevelNames = map[slog.Leveler]string{
	LevelNotice: "NOC",
	LevelFatal:  "FAL",
}

type implement struct {
	log *slog.Logger
}

type noCfg struct{}

var implementer *implement

type InalogHandler struct {
	slog.Handler
	skipCaller int
}

func createJsonHandler(cfg Cfg) slog.Handler {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
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

func Init(cfg Cfg) {
	var handler slog.Handler

	if cfg.CustomFunc != nil {
		handler = cfg.CustomFunc(cfg)
	} else if cfg.TextLog {
		handler = createTextHandler(cfg)
	} else {
		handler = createJsonHandler(cfg)
	}
	h := &InalogHandler{handler, 3}
	logger := slog.New(h)

	slog.SetDefault(logger)
	implementer = &implement{
		log: logger,
	}
}

func Log() Interface {
	return &noCfg{}
}

func (i *noCfg) Debug(msg string, attr ...any) {
	newWriter(writer{
		logImpl: serviceContext(implementer.log),
	}).logHandling(slog.LevelDebug, msg, attr...)
}

func (i *noCfg) Info(msg string, attr ...any) {
	newWriter(writer{
		logImpl: serviceContext(implementer.log),
	}).logHandling(slog.LevelInfo, msg, attr...)
}

func (i *noCfg) Notice(msg string, attr ...any) {
	newWriter(writer{
		logImpl: serviceContext(implementer.log),
	}).logHandling(LevelNotice, msg, attr...)
}

func (i *noCfg) Warn(msg string, attr ...any) {
	newWriter(writer{
		logImpl: serviceContext(implementer.log),
	}).logHandling(slog.LevelWarn, msg, attr...)
}

func (i *noCfg) Error(msg string, attr ...any) {
	newWriter(writer{
		logImpl: serviceContext(implementer.log),
	}).logHandling(slog.LevelError, msg, attr...)
}

func (i *noCfg) Fatal(msg string, attr ...any) {
	newWriter(writer{
		logImpl: serviceContext(implementer.log),
	}).logHandling(LevelFatal, msg, attr...)
}
