package inalog

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type FiberHTTPLogParam struct {
	FiberCtx  *fiber.Ctx
	StartTime time.Time
}

func FiberCtxDeviceBuilder(f *fiber.Ctx) map[string]interface{} {
	appVersion := f.Get("X-App-Version")
	deviceType := f.Get("X-Device-Type")

	if deviceType == "" {
		deviceType = "Web"
	}

	if appVersion != "" {
		deviceType = "Mobile"
	}

	return map[string]interface{}{
		"ID":         f.Get("X-Device-ID"),
		"Info":       f.Get("X-Device-Info"),
		"AppVersion": appVersion,
		"Type":       deviceType,
	}
}

func FiberCtxHttpBuilder(f *fiber.Ctx) map[string]interface{} {
	return map[string]interface{}{
		"method":          f.Method(),
		"url":             string(f.Request().URI().Path()),
		"requestSize":     f.Request().Header.ContentLength(),
		"status":          f.Response().StatusCode(),
		"userAgent":       string(f.Request().Header.UserAgent()),
		"sourceIp":        f.IP(),
		"sourceIps":       f.IPs(),
		"referer":         string(f.Request().Header.Referer()),
		"hostname":        f.Hostname(),
		"X-Request-ID":    f.Get("X-Request-ID"),
		"X-Forwarded-For": f.Get("X-Forwarded-For"),
		"X-Real-IP":       f.Get("X-Real-IP"),
	}
}

func cleanJSONPrint(input string) string {
	return strings.NewReplacer("\n", "", "\t", "", "\\", "", "\"", "'").Replace(input)
}

func FiberHTTPLog(param FiberHTTPLogParam) {
	printPayload, _ := strconv.ParseBool(os.Getenv("INALOG_PRINT_PAYLOAD"))
	printAccess, _ := strconv.ParseBool(os.Getenv("INALOG_ACCESS_LOG"))
	printError, _ := strconv.ParseBool(os.Getenv("INALOG_ERROR_LOG"))

	minStatusToPrint := int(100)

	fiberCtx := param.FiberCtx
	statusCode := fiberCtx.Response().StatusCode()

	queries := fiberCtx.Queries()
	enforced := false
	forced, ok := queries["_InalogForcePrint"]
	if ok {
		enforced, _ = strconv.ParseBool(forced)
	}

	printBody := (statusCode >= minStatusToPrint && printPayload) || enforced
	printHeaders := (statusCode >= minStatusToPrint && printPayload) || enforced
	printQuery := (statusCode >= minStatusToPrint && printPayload) || enforced

	data := FiberCtxHttpBuilder(fiberCtx)
	getDuration := time.Since(param.StartTime)
	data["duration"] = getDuration.String()
	data["durationInMs"] = getDuration.Milliseconds()

	if printBody {
		data["req_body"] = cleanJSONPrint(string(fiberCtx.BodyRaw()))
	}

	if printHeaders {
		reqHeaders, _ := json.Marshal(fiberCtx.GetReqHeaders())
		data["headers"] = cleanJSONPrint(string(reqHeaders))
	}

	if printQuery {
		queryStrings, _ := json.Marshal(fiberCtx.Queries())
		data["query_params"] = cleanJSONPrint(string(queryStrings))
	}

	ctx := context.WithValue(WithFiberCtx(fiberCtx.Context()), CtxKeyHttp, data)
	if printAccess && statusCode >= 200 && statusCode < 300 {
		LogWith(WithCfg{Ctx: ctx, Skip: 1}).
			Notice(fmt.Sprintf(
				"%s %s",
				fiberCtx.Request().Header.Method(),
				string(fiberCtx.Request().URI().Path()),
			))
	}

	if printError && statusCode >= 300 {
		LogWith(WithCfg{Ctx: ctx, Skip: 1}).
			Notice(fmt.Sprintf(
				"%s %s",
				fiberCtx.Request().Header.Method(),
				string(fiberCtx.Request().URI().Path()),
			))
	}
}

func FiberInheriCtx(f *fiber.Ctx) context.Context {
	ctx := context.Background()

	for _, v := range CtxList {
		ctxVal := f.Context().Value(v)

		if ctxVal != nil {
			ctx = context.WithValue(ctx, v, ctxVal)
		}
	}

	return ctx
}

func NewFiberMiddleware() fiber.Handler {
	return func(f *fiber.Ctx) error {
		startTime := time.Now()
		err := f.Next()
		FiberHTTPLog(FiberHTTPLogParam{
			f,
			startTime,
		})
		return err
	}
}

func HttpHeaderToSlog(header http.Header) slog.Attr {
	var headers []any
	for key, values := range header {
		if len(values) == 0 {
			headers = append(headers, slog.Any(key, nil))
		} else if len(values) == 1 {
			headers = append(headers, slog.String(key, values[0]))
		} else {
			headers = append(headers, slog.Any(key, values))
		}
	}
	return slog.Group("header", headers...)
}
