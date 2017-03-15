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

type AudioReceiver struct {
	*Transport
	sync.Mutex
	Name             string
	Device           string
	PulsePlaybackIdx int32
	samplerate       int32
	channels         int32
}

func NewAudioReceiver(name string, device string, samplerate int32, channels int32) *AudioReceiver {
	idx := alsa.Pa_init(name, device, samplerate, channels, 1)
	if idx < 0 {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Failed to initialize pulseaudio: %v", idx))
		return nil
	}

	ar := &AudioReceiver{}
	ar.Transport = &Transport{}
	ar.Name = name
	ar.Device = device
	ar.PulsePlaybackIdx = idx
	ar.initialize(samplerate, channels)
	return ar
}

func (ar *AudioReceiver) initialize(samplerate int32, channels int32) {
	ar.samplerate = samplerate
	ar.channels = channels
}

func (ar *AudioReceiver) BeginReception(dataCallback func(b *[]byte)) (err error) {
	if ar.Conn == nil {
		err = errors.New("Cannot begin reception before connection to transmitter is established")
		return
	}
	data := make([]byte, ar.samplerate)
	for {
		ar.Lock()
		if _, err = ar.Read(data); err != nil {
			ar.Unlock()
			return
		}
		ar.Unlock()
		if dataCallback != nil {
			dataCallback(&data)
		}
		if ret := alsa.Pa_handle_write(ar.PulsePlaybackIdx, &data, int32(len(data))); ret != 0 {
			err = errors.New(fmt.Sprintf("Failed to write data to pulseaudio: %v", ret))
			return
		}
	}
	return
}

func (ar *AudioReceiver) Listen(proto string, addr string) (err error) {
	var listener net.Listener
	var kcpListener *kcp.Listener

	fmt.Println("Listening for a connection...")
	ar.Lock()
	defer ar.Unlock()
	switch proto {
	case "tcp":
		if listener, err = net.Listen("tcp", addr); err != nil {
			return
		}
		ar.Conn, err = listener.Accept()
	case "udp":
		if kcpListener, err = kcp.ListenWithOptions(addr, nil, 10, 3); err != nil {
			return
		} else {
			ar.Conn, err = kcpListener.AcceptKCP()
		}
	}
	fmt.Println("Received connection from:", ar.Conn.RemoteAddr())
	return
}
