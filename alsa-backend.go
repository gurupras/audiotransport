package audiotransport

import (
	"errors"
	"fmt"

	"github.com/gurupras/audiotransport/alsa"
)

type AlsaBackend struct {
	*Backend
}

func (ab *AlsaBackend) Init(name, device string, samplerate, channels uint32, isPlayback bool) (err error) {
	if ab.Backend == nil {
		ab.Backend = &Backend{}
	}
	ab.Backend.Init(name, device, samplerate, channels)

	var playback int32 = 0
	if isPlayback {
		playback = 1
	}

	idx := alsa.Alsa_init(device, int32(samplerate), int32(channels), playback)
	if idx < 0 {
		err = errors.New(fmt.Sprintf("Failed to initialize ALSA device: %v", idx))
	}
	return
}

func (ab *AlsaBackend) Read(buf []byte, len uint32) (int, error) {
	return int(alsa.Alsa_readi(ab.HandleIdx, &buf, int32(len))), nil
}

func (ab *AlsaBackend) Write(buf []byte, len uint32) (int, error) {
	return int(alsa.Alsa_writei(ab.HandleIdx, &buf, int32(len))), nil
}
