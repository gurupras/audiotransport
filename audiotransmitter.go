package audiotransport

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/gurupras/audiotransport/alsa"
)

type AudioTransmitter struct {
	*UdpClient
	sync.Mutex
	Name            string
	Device          string
	PulseCaptureIdx int32
	samplerate      int32
	channels        int32
}

func NewAudioTransmitter(name string, device string, samplerate int32, channels int32) *AudioTransmitter {
	idx := alsa.Pa_init(name, device, samplerate, channels, 0)
	if idx < 0 {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Failed to initialize pulseaudio: %v", idx))
		return nil
	}

	at := &AudioTransmitter{}
	at.Name = name
	at.Device = device
	at.UdpClient = NewUDPClient()
	at.PulseCaptureIdx = idx
	at.initialize(samplerate, channels)
	return at
}

func (at *AudioTransmitter) initialize(samplerate int32, channels int32) {
	at.samplerate = samplerate
	at.channels = channels
}

func (at *AudioTransmitter) BeginTransmission() (err error) {
	if at.UdpClient.Conn == nil {
		err = errors.New("Cannot begin transmission before connection to receiver is established")
		return
	}
	size := at.samplerate
	buf := make([]byte, size)
	for {
		alsa.Pa_handle_read(at.PulseCaptureIdx, &buf, size)
		at.Lock()
		if _, err = at.WriteBytes(buf); err != nil {
			err = errors.New(fmt.Sprintf("Failed to send data over transport: %v", err))
			return
		}
		at.Unlock()
	}
	return
}

func (at *AudioTransmitter) Connect(addr string) (err error) {
	at.Lock()
	err = at.UdpClient.Connect(addr)
	at.Unlock()
	return
}
