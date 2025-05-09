package inalog

import (
	"context"
	"encoding/json"
	"fmt"
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
		"referer":         f.Request().Header.Referer(),
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
	fiberCtx := param.FiberCtx
	data := FiberCtxHttpBuilder(fiberCtx)
	getDuration := time.Since(param.StartTime)
	data["duration"] = getDuration.String()
	data["durationInMs"] = getDuration.Milliseconds()

	if fiberCtx.Response().StatusCode() >= 100 {
		queryStrings, _ := json.Marshal(fiberCtx.Queries())
		reqHeaders, _ := json.Marshal(fiberCtx.GetReqHeaders())
		data["query_params"] = cleanJSONPrint(string(queryStrings))
		data["headers"] = cleanJSONPrint(string(reqHeaders))
		data["req_body"] = cleanJSONPrint(string(fiberCtx.BodyRaw()))
	}

	print, _ := strconv.ParseBool(os.Getenv("INALOG_ACCESS_LOG"))
	if print {
		LogWith(WithCfg{Ctx: context.WithValue(fiberCtx.Context(), CtxKeyHttp, data), Skip: 1}).
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
