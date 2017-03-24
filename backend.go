package audiotransport

import (
	"errors"
	"io"
)

type BackendInterface interface {
	Init(name, device string, samplerate, channels uint32, isPlayback bool) error
	io.Reader
	io.Writer
	GetBufferSize() uint32
	GetLatency() (int64, error)
}

type Backend struct {
	Name       string
	Device     string
	HandleIdx  int32
	SampleRate uint32
	Channels   uint32
}

func (b *Backend) Init(name, device string, samplerate, channels uint32) {
	b.Name = name
	b.Device = device
	b.SampleRate = samplerate
	b.Channels = channels
}

func (b *Backend) GetBufferSize() uint32 {
	return 512
}

func (b *Backend) GetLatency() (int64, error) {
	return -1, errors.New("Unimplemented")
}
