package inalog

type Interface interface {
	Debug(string, ...any)
	Info(string, ...any)
	Notice(string, ...any)
	Warn(string, ...any)
	Error(string, ...any)
	Fatal(string, ...any)
}
