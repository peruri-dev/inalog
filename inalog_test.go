package inalog

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/peruri-dev/inalog/testhelper"
	"github.com/stretchr/testify/assert"
)

func TestInalog(t *testing.T) {
	assert := assert.New(t)

	now := time.Now()
	output, _ := testhelper.CaptureOutput(func() error {
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

	var jsonOutput testhelper.ExpectedJSON
	err := json.Unmarshal(output, &jsonOutput)
	assert.NoError(err)
	assert.WithinDuration(now, jsonOutput.Time, 5*time.Second)
	assert.Equal("INFO", jsonOutput.Level)
	assert.Equal("this is informative", jsonOutput.Msg)
	assert.Contains(jsonOutput.Source.File, "/inalog_test.go")
	assert.Contains(jsonOutput.Source.Function, "/inalog.TestInalog.")
	assert.Equal(24, jsonOutput.Source.Line)

	assert.Equal("example-app", jsonOutput.Service.Name)
	assert.Equal("testing", jsonOutput.Service.Env)
	assert.Equal("v0.0.0-beta", jsonOutput.Service.Version)
	assert.NotZero(jsonOutput.Service.Pid)
}

func TestInalogCfgMessageKey(t *testing.T) {
	assert := assert.New(t)

	output, _ := testhelper.CaptureOutput(func() error {
		Init(Cfg{
			MessageKey: true,
		})
		Log().Warn("this is warning")
		return nil
	})

	assert.NotEmpty(output)
	var jsonOutput testhelper.ExpectedJSON
	err := json.Unmarshal(output, &jsonOutput)
	assert.NoError(err)

	assert.Equal("WARN", jsonOutput.Level)
	assert.Equal("this is warning", jsonOutput.Message)
}
