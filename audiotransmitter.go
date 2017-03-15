package audiotransport

import (
	"errors"
	"fmt"
	"net"
	"os"
	"sync"

	"github.com/gurupras/audiotransport/alsa"
	"github.com/xtaci/kcp-go"
)

type AudioTransmitter struct {
	*Transport
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
	at.Transport = &Transport{}
	at.Name = name
	at.Device = device
	at.PulseCaptureIdx = idx
	at.initialize(samplerate, channels)
	return at
}

func (at *AudioTransmitter) initialize(samplerate int32, channels int32) {
	at.samplerate = samplerate
	at.channels = channels
}

func (at *AudioTransmitter) BeginTransmission() (err error) {
	if at.Conn == nil {
		err = errors.New("Cannot begin transmission before connection to receiver is established")
		return
	}
	size := at.samplerate / 2
	buf := make([]byte, size)
	for {
		alsa.Pa_handle_read(at.PulseCaptureIdx, &buf, size)
		at.Lock()
		fmt.Printf("Attempting to send %d bytes\n", len(buf))
		if _, err = at.Write(buf); err != nil {
			err = errors.New(fmt.Sprintf("Failed to send data over transport: %v", err))
			return
		}
		at.Unlock()
		fmt.Printf("Sent %d bytes\n", len(buf))
	}
	return
}

func (at *AudioTransmitter) Connect(proto string, addr string) (err error) {
	at.Lock()
	defer at.Unlock()
	switch proto {
	case "tcp":
		at.Conn, err = net.Dial("tcp", addr)
	case "udp":
		at.Conn, err = kcp.DialWithOptions(addr, nil, 10, 3)
	}
	return
}
