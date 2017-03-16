package audiotransport

import (
	"errors"
	"fmt"
	"net"
	"os"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/gurupras/audiotransport/alsa"
	"github.com/xtaci/kcp-go"
)

type AudioReceiver struct {
	Transport
	ApiType
	sync.Mutex
	Name        string
	Device      string
	PlaybackIdx int32
	samplerate  int32
	channels    int32
}

func NewAudioReceiver(apiType ApiType, name string, device string, samplerate int32, channels int32) *AudioReceiver {
	var idx int32
	switch apiType {
	case ALSA_API:
		idx = alsa.Alsa_init(device, samplerate, channels, 1)
	case PULSE_API:
		idx = alsa.Pa_init(name, device, samplerate, channels, 1)
	}

	if idx < 0 {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Failed to initialize %s: %v", apiType.ApiString(), idx))
		return nil
	}

	ar := &AudioReceiver{}
	ar.ApiType = apiType
	ar.Name = name
	ar.Device = device
	ar.PlaybackIdx = idx
	ar.initialize(samplerate, channels)
	return ar
}

func (ar *AudioReceiver) initialize(samplerate int32, channels int32) {
	ar.samplerate = samplerate
	ar.channels = channels
}

func (ar *AudioReceiver) BeginReception(dataCallback func(b *[]byte)) (err error) {
	if ar.Transport == nil {
		err = errors.New("Cannot begin reception before connection to transmitter is established")
		return
	}
	bufsize := ar.GetBufferSize(ar.samplerate, ar.channels)
	data := make([]byte, bufsize)
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
		var ret int32
		switch ar.ApiType {
		case ALSA_API:
			ret = alsa.Alsa_writei(ar.PlaybackIdx, &data, int32(len(data)))
		case PULSE_API:
			ret = alsa.Pa_handle_write(ar.PlaybackIdx, &data, int32(len(data)))
		}
		if ret != 0 {
			err = errors.New(fmt.Sprintf("Failed to write data to %s: %v", ar.ApiType.ApiString(), ret))
			return
		}
	}
	return
}

func (ar *AudioReceiver) Listen(proto string, addr string) (err error) {
	var listener net.Listener
	var kcpListener *kcp.Listener
	var conn net.Conn

	log.Info("Listening for a connection...")
	ar.Lock()
	defer ar.Unlock()
	switch proto {
	case "tcp":
		if listener, err = net.Listen("tcp", addr); err != nil {
			return
		}
		conn, err = listener.Accept()
		transport := &BaseTransport{}
		transport.Conn = conn
		ar.Transport = transport

	case "kcp":
		if kcpListener, err = kcp.ListenWithOptions(addr, nil, 10, 3); err != nil {
			return
		} else {
			conn, err = kcpListener.AcceptKCP()
		}
		transport := &BaseTransport{}
		transport.Conn = conn
		ar.Transport = transport
	case "udp":
		server := NewUDPServer()
		wg := sync.WaitGroup{}
		wg.Add(1)
		callback := func(transport Transport) {
			defer wg.Done()
			ar.Transport = transport
		}
		if err = server.Listen(addr, callback); err != nil {
			return
		}
		wg.Wait()
	}
	log.Info("Received connection from:", ar.Transport.RemoteAddr())
	return
}
