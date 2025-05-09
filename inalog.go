package inalog

import (
	"log/slog"
	"os"
)

type Cfg struct {
	Source bool
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

type implement struct {
	log *slog.Logger
}

type noCfg struct{}

var implementer *implement

func Init(cfg Cfg) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
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

			return a
		},
	}))

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
