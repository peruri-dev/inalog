package inalog

import (
	"encoding/json"
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func captureOutput(f func() error) ([]byte, error) {
	orig := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	err := f()
	os.Stdout = orig
	w.Close()
	out, _ := io.ReadAll(r)
	return out, err
}

type ExpectedSourceJSON struct {
	Function string `json:"function"`
	File     string `json:"file"`
	Line     int    `json:"line"`
}

type ExpectedServiceJSON struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Env     string `json:"env"`
	Pid     int    `json:"pid"`
}

type ExpectedJSON struct {
	Time      time.Time           `json:"time"`
	Level     string              `json:"level"`
	Source    ExpectedSourceJSON  `json:"source"`
	Msg       string              `json:"msg"`
	Message   string              `json:"message"`
	Service   ExpectedServiceJSON `json:"service"`
	HTTP      ExpectedHTTPJSON    `json:"http"`
	RequestId string              `json:"requestId"`
	TraceId   string              `json:"traceId"`
	SpanId    string              `json:"spanId"`
}

func TestInalog(t *testing.T) {
	assert := assert.New(t)

	now := time.Now()
	output, _ := captureOutput(func() error {
		t.Setenv("INALOG_SERVICE_VERSION", "v0.0.0-beta")
		t.Setenv("INALOG_SERVICE_NAME", "example-app")
		t.Setenv("INALOG_SERVICE_ENV", "testing")

		Init(Cfg{
			Source: true,
		})
		Log().Info("this is informative")
		return nil
	})

	assert.NotEmpty(output)

	var jsonOutput ExpectedJSON
	err := json.Unmarshal(output, &jsonOutput)
	assert.NoError(err)
	assert.WithinDuration(now, jsonOutput.Time, 5*time.Second)
	assert.Equal("INFO", jsonOutput.Level)
	assert.Equal("this is informative", jsonOutput.Msg)
	assert.Contains(jsonOutput.Source.File, "/inalog_test.go")
	assert.Contains(jsonOutput.Source.Function, "/inalog.TestInalog.")
	assert.Equal(62, jsonOutput.Source.Line)

	assert.Equal("example-app", jsonOutput.Service.Name)
	assert.Equal("testing", jsonOutput.Service.Env)
	assert.Equal("v0.0.0-beta", jsonOutput.Service.Version)
	assert.NotZero(jsonOutput.Service.Pid)
}

func TestInalogCfgMessageKey(t *testing.T) {
	assert := assert.New(t)

	output, _ := captureOutput(func() error {
		Init(Cfg{
			MessageKey: true,
		})
		Log().Warn("this is warning")
		return nil
	})

	assert.NotEmpty(output)
	var jsonOutput ExpectedJSON
	err := json.Unmarshal(output, &jsonOutput)
	assert.NoError(err)

	assert.Equal("WARN", jsonOutput.Level)
	assert.Equal("this is warning", jsonOutput.Message)
}
