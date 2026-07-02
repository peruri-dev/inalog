package testhelper

import "time"

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
	RequestId string              `json:"requestId"`
	TraceId   string              `json:"traceId"`
	SpanId    string              `json:"spanId"`
}
