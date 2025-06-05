package inalog

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

type ExpectedHTTPJSON struct {
	Method      string   `json:"method"`
	URL         string   `json:"url"`
	RequestSize int      `json:"requestSize"`
	Status      int      `json:"status"`
	UserAgent   string   `json:"userAgent"`
	SourceIP    string   `json:"sourceIp"`
	SourceIPs   []string `json:"sourceIps"`
	Referer     string   `json:"referer"`
	Hostname    string   `json:"hostname"`
	ForwaredFor string   `json:"X-Forwarded-For"`
	RealIP      string   `json:"X-Real-IP"`
}

func TestFiberHTTPLog(t *testing.T) {
	assert := assert.New(t)
	t.Setenv("INALOG_ACCESS_LOG", "true")

	startTime := time.Now()

	app := fiber.New()
	app.Get("/", func(f *fiber.Ctx) error {
		output, _ := captureOutput(func() error {
			f.Context().SetUserValue(CtxKeyRequestID, "abcd-wkwkwk")
			f.Context().SetUserValue(CtxKeyTraceID, "efgh-wkwk")
			f.Context().SetUserValue(CtxKeySpanID, "ijkl-wk")

			Init(Cfg{})

			FiberHTTPLog(FiberHTTPLogParam{
				f,
				startTime,
			})
			return nil
		})

		return f.SendString(string(output))
	})

	req := httptest.NewRequest(fiber.MethodGet, "/", nil)
	req.Host = "inalog.com"
	req.Header.Add("X-Real-IP", "1.1.1.1")
	req.Header.Add("Referer", "http://gogel.com")

	resp, err := app.Test(req)
	assert.NoError(err)

	output, err := io.ReadAll(resp.Body)
	assert.NoError(err)

	var jsonOutput ExpectedJSON
	err = json.Unmarshal(output, &jsonOutput)
	assert.NoError(err)

	assert.Equal("GET", jsonOutput.HTTP.Method)
	assert.Equal("/", jsonOutput.HTTP.URL)
	assert.Equal(200, jsonOutput.HTTP.Status)
	assert.Equal("inalog.com", jsonOutput.HTTP.Hostname)
	assert.Equal("1.1.1.1", jsonOutput.HTTP.RealIP)
	assert.Equal("abcd-wkwkwk", jsonOutput.RequestId)
	assert.Equal("efgh-wkwk", jsonOutput.TraceId)
	assert.Equal("ijkl-wk", jsonOutput.SpanId)
	assert.Equal("http://gogel.com", jsonOutput.HTTP.Referer)
}
