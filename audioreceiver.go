package audiotransport

import (
	"errors"
	"fmt"
	"net"
	"os"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/xtaci/kcp-go"
)

type AudioReceiver struct {
	Transport
	ApiType
	Backend BackendInterface
	sync.Mutex
}

func NewAudioReceiver(apiType ApiType, name string, device string, samplerate int32, channels int32) *AudioReceiver {
	var backend BackendInterface

	switch apiType {
	case ALSA_API:
		backend = &AlsaBackend{}
	case PULSE_API:
		backend = &PulseBackend{}
	}

	if err := backend.Init(name, device, samplerate, channels, 1); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Failed to initialize %s: %v", apiType.ApiString(), err))
		return nil
	}

	ar := &AudioReceiver{}
	ar.ApiType = apiType
	ar.Backend = backend
	return ar
}

func (ar *AudioReceiver) BeginReception(dataCallback func(b *[]byte)) (err error) {
	if ar.Transport == nil {
		err = errors.New("Cannot begin reception before connection to transmitter is established")
		return
	}
	bufsize := ar.Backend.GetBufferSize()
	buf := make(SoundBytes, bufsize)
	bufBytes := []byte(buf)
	for {
		ar.Lock()
		if _, err = ar.Read(bufBytes); err != nil {
			ar.Unlock()
			return
		}
		ar.Unlock()
		if dataCallback != nil {
			dataCallback(&bufBytes)
		}
		var ret int32
		if ret = ar.Backend.Write(bufBytes, bufsize); ret != 0 {
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
