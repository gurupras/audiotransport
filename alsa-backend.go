package audiotransport

import (
	"errors"
	"fmt"

	"github.com/gurupras/audiotransport/alsa"
)

type AlsaBackend struct {
	*Backend
}

func (ab *AlsaBackend) Init(name, device string, samplerate, channels, isPlayback int32) (err error) {
	if ab.Backend == nil {
		ab.Backend = &Backend{}
	}
	ab.Backend.Init(name, device, samplerate, channels)

	idx := alsa.Alsa_init(device, samplerate, channels, isPlayback)
	if idx < 0 {
		err = errors.New(fmt.Sprintf("Failed to initialize ALSA device: %v", idx))
	}
	return
}

func (ab *AlsaBackend) Read(buf []byte, len int32) int32 {
	return alsa.Alsa_readi(ab.HandleIdx, &buf, len)
}

func (ab *AlsaBackend) Write(buf []byte, len int32) int32 {
	return alsa.Alsa_writei(ab.HandleIdx, &buf, len)
}
