package mic

import (
	"io"

	"github.com/cocoonlife/goalsa"
)

// NewReader returns an interface to the microphone
// which acts just like an io.Reader.
func NewReader() io.Reader {
	return reader{}
}

type reader struct{}

func (reader) Read(p []byte) (n int, err error) {
	dev, err := alsa.NewCaptureDevice("default", 1, alsa.FormatU8, 8000, alsa.BufferParams{})
	if err != nil {
		return
	}
	b := make([]int8, len(p))
	n, err = dev.Read(b)
	for i := 0; i < n; i++ {
		p[i] = byte(b[i])
	}
	return
}
