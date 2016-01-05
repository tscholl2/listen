package mic

import (
	"fmt"
	"os"
	"testing"

	"github.com/cocoonlife/goalsa"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMic(t *testing.T) {
	dev, err := alsa.NewCaptureDevice("default", 1, alsa.FormatU8, 8000, alsa.BufferParams{})
	if err != nil {
		t.Errorf("Err initializing mic:\n\t%s", err)
		t.FailNow()
	}
	d := 5 //seconds to record
	r := make([]int8, 8000*d)
	fmt.Printf("Recording for %d seconds...\n", d)
	n, err := dev.Read(r)
	require.Nil(t, err)
	require.Equal(t, n, len(r))
	//cast to bytes
	b := make([]byte, len(r))
	for i := 0; i < len(b); i++ {
		b[i] = byte(r[i])
	}
	fmt.Printf("Done recording. Saving file.")
	f, err := os.OpenFile("./a.wav", os.O_CREATE|os.O_WRONLY, 0666)
	require.Nil(t, err)
	_, err = f.Write(b)
	require.Nil(t, err)
	err = f.Close()
	require.Nil(t, err)
}

func TestReader(t *testing.T) {
	r := NewReader()
	b := make([]byte, 1000)
	n, err := r.Read(b)
	assert.Nil(t, err)
	assert.Equal(t, n, len(b))
	s := 0
	for _, x := range b {
		s += int(x)
	}
	assert.NotEqual(t, s, 0)
}
