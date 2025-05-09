package inalog

import (
	"context"
	"log/slog"
)

type WithCfg struct {
	Ctx  context.Context
	Skip int
}

func LogWith(wcfg WithCfg) Interface {
	if implementer == nil {
		Init(Cfg{})
	}

	return &wcfg
}

func (i WithCfg) Debug(msg string, attr ...any) {
	newWriter(writer{
		logImpl:    serviceContext(implementer.log),
		ctx:        i.Ctx,
		skipCaller: i.Skip,
	}).logHandling(slog.LevelDebug, msg, attr...)
}

func (i WithCfg) Info(msg string, attr ...any) {
	newWriter(writer{
		logImpl:    serviceContext(implementer.log),
		ctx:        i.Ctx,
		skipCaller: i.Skip,
	}).logHandling(slog.LevelInfo, msg, attr...)
}

func (i WithCfg) Notice(msg string, attr ...any) {
	newWriter(writer{
		logImpl:    serviceContext(implementer.log),
		ctx:        i.Ctx,
		skipCaller: i.Skip,
	}).logHandling(LevelNotice, msg, attr...)
}

func (i WithCfg) Warn(msg string, attr ...any) {
	newWriter(writer{
		logImpl:    serviceContext(implementer.log),
		ctx:        i.Ctx,
		skipCaller: i.Skip,
	}).logHandling(slog.LevelWarn, msg, attr...)
}

func (i WithCfg) Error(msg string, attr ...any) {
	newWriter(writer{
		logImpl:    serviceContext(implementer.log),
		ctx:        i.Ctx,
		skipCaller: i.Skip,
	}).logHandling(slog.LevelError, msg, attr...)
}

func (i WithCfg) Fatal(msg string, attr ...any) {
	newWriter(writer{
		logImpl:    serviceContext(implementer.log),
		ctx:        i.Ctx,
		skipCaller: i.Skip,
	}).logHandling(LevelFatal, msg, attr...)
}
