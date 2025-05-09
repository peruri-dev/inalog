package inalog

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type ExternalHTTPCallLogParam struct {
	Ctx        context.Context
	Response   *http.Response
	ReqBodyRaw []byte
	StartTime  time.Time
}

func HTTPRequestCtxBuilder(ctx context.Context, res *http.Response) map[string]interface{} {
	return map[string]interface{}{
		"method":          res.Request.Method,
		"url":             res.Request.URL.String(),
		"requestSize":     res.Request.ContentLength,
		"status":          res.StatusCode,
		"userAgent":       res.Request.UserAgent(),
		"sourceIp":        res.Request.RemoteAddr,
		"referer":         res.Request.Referer(),
		"hostname":        res.Request.URL.Hostname(),
		"X-Forwarded-For": res.Request.Header.Get("X-Forwarded-For"),
		"X-Real-IP":       res.Request.Header.Get("X-Real-IP"),
	}
}

func ExternalHTTPCallLog(param ExternalHTTPCallLogParam) {
	data := HTTPRequestCtxBuilder(param.Ctx, param.Response)
	getDuration := time.Since(param.StartTime)
	data["duration"] = getDuration.String()
	data["durationInMs"] = getDuration.Milliseconds()

	if param.Response.StatusCode >= 300 {
		queryStrings, _ := json.Marshal(param.Response.Request.URL.Query())
		reqHeaders, _ := json.Marshal(param.Response.Request.Header)
		data["query_params"] = string(queryStrings)
		data["headers"] = string(reqHeaders)
		data["req_body"] = strings.NewReplacer("\n", "", "\t", "", "\\", "", " ", "").Replace(string(param.ReqBodyRaw))
	}

	LogWith(WithCfg{Ctx: context.WithValue(param.Ctx, CtxKeyHttpRequest, data), Skip: 1}).
		Notice(fmt.Sprintf(
			"External HTTP Request %s %s",
			param.Response.Request.Method,
			param.Response.Request.URL.Path,
		))
}
