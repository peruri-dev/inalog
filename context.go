package inalog

import (
	"context"
	"log/slog"
	"os"
)

type CtxKey string

var CtxKeyErrContext CtxKey = "errorContext"
var CtxKeyPayloadContext CtxKey = "payloadContext"
var CtxKeyAuditContext CtxKey = "auditContext"
var CtxKeyHttp CtxKey = "http"
var CtxKeyHttpRequest CtxKey = "httpRequest"
var CtxKeyRequestID CtxKey = "requestId"
var CtxKeyDevice CtxKey = "deviceContext"

// Put context here, to make it printable in log
var CtxList = []CtxKey{
	CtxKeyRequestID,
	CtxKeyHttp,
	CtxKeyHttpRequest,
	CtxKeyErrContext,
	CtxKeyPayloadContext,
	CtxKeyDevice,
	"spanId",
	"traceId",
	"span.id",
	"trace.id",
	"transaction.id",
}

func contextParser(ctx context.Context) []slog.Attr {
	ctxDecorator := []slog.Attr{}

	for _, v := range CtxList {
		ctxVal := ctx.Value(v)

		if ctxVal != nil {
			ctxDecorator = append(ctxDecorator, slog.Any(string(v), ctxVal))
		}
	}

	return ctxDecorator
}

func serviceContext(l *slog.Logger) *slog.Logger {
	return l.With(
		slog.Group("service",
			slog.String("name", os.Getenv("INALOG_SERVICE_NAME")),
			slog.String("version", os.Getenv("INALOG_SERVICE_VERSION")),
			slog.String("env", os.Getenv("INALOG_SERVICE_ENV")),
			slog.Int("pid", os.Getpid()),
		),
	)
}

func ErrorCtx(err error) slog.Attr {
	return slog.Any(string(CtxKeyErrContext), err)
}

func PayloadCtx(value interface{}) slog.Attr {
	return slog.Any(string(CtxKeyPayloadContext), value)
}

func AuditCtx(value interface{}) slog.Attr {
	return slog.Any(string(CtxKeyAuditContext), value)
}
