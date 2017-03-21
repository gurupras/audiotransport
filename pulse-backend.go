package audiotransport

import (
	"errors"
	"fmt"

	"github.com/gurupras/audiotransport/alsa"
)

type PulseBackend struct {
	*Backend
}

func (pb *PulseBackend) Init(name, device string, samplerate, channels, isPlayback int32) (err error) {
	if pb.Backend == nil {
		pb.Backend = &Backend{}
	}
	pb.Backend.Init(name, device, samplerate, channels)

	idx := alsa.Pa_init(name, device, samplerate, channels, isPlayback)
	if idx < 0 {
		err = errors.New(fmt.Sprintf("Failed to initialize PULSE: %v", idx))
	}
	return
}

func (pb *PulseBackend) Read(buf []byte, len int32) int32 {
	return alsa.Pa_handle_read(pb.HandleIdx, &buf, len)
}

func (pb *PulseBackend) Write(buf []byte, len int32) int32 {
	return alsa.Pa_handle_write(pb.HandleIdx, &buf, len)
}

func (pb *PulseBackend) GetLatency() int32 {
	return alsa.Pa_get_latency(pb.HandleIdx)
}
