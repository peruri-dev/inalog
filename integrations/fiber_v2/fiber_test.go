package inalog

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/peruri-dev/inalog"
	"github.com/peruri-dev/inalog/testhelper"
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

type ExpectedJSON struct {
	testhelper.ExpectedJSON
	HTTP ExpectedHTTPJSON `json:"http"`
}

type ResultJSON struct {
	Success bool `json:"success"`
}

func TestFiberHTTPLog(t *testing.T) {
	assert := assert.New(t)
	t.Setenv("INALOG_ACCESS_LOG", "true")

	//startTime := time.Now()

	app := fiber.New()
	app.Use(NewFiberMiddleware())

	app.Get("/", func(f *fiber.Ctx) error {
		inalog.Init(inalog.Cfg{})

		f.Context().SetUserValue(inalog.CtxKeyRequestID, "abcd-wkwkwk")
		f.Context().SetUserValue(inalog.CtxKeyTraceID, "efgh-wkwk")
		f.Context().SetUserValue(inalog.CtxKeySpanID, "ijkl-wk")

		return f.SendString(`OK`)
	})

	var resp *http.Response
	output, err := testhelper.CaptureOutput(func() error {
		req := httptest.NewRequest(fiber.MethodGet, "/", nil)
		req.Host = "inalog.com"
		req.Header.Add("X-Real-IP", "1.1.1.1")
		req.Header.Add("Referer", "http://gogel.com")

		var err error

		resp, err = app.Test(req)
		return err
	})

	assert.NoError(err)

	respOutput, err := io.ReadAll(resp.Body)
	assert.NoError(err)
	assert.Equal([]byte("OK"), respOutput)

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
