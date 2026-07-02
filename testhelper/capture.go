package testhelper

import (
	"io"
	"os"
)

func CaptureOutput(f func() error) ([]byte, error) {
	orig := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	err := f()
	os.Stdout = orig
	w.Close()
	out, _ := io.ReadAll(r)
	return out, err
}
