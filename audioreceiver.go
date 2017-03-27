package audiotransport

import (
	"errors"
	"fmt"
	"os"
	"sync"

	log "github.com/Sirupsen/logrus"
)

type AudioReceiver struct {
	Transport
	ApiType
	Backend           BackendInterface
	ReceptionCallback func(data *[]byte) (err error)
	sync.Mutex
}

func NewAudioReceiver(apiType ApiType, name string, device string, samplerate, channels uint32) *AudioReceiver {
	var backend BackendInterface

	switch apiType {
	case ALSA_API:
		backend = &AlsaBackend{}
	case PULSE_API:
		backend = &PulseSimpleBackend{}
	case FILE_API:
		backend = &FileBackend{}
	}

	if err := backend.Init(name, device, samplerate, channels, true); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Failed to initialize %s: %v", apiType.ApiString(), err))
		return nil
	}

	ar := &AudioReceiver{}
	ar.ApiType = apiType
	ar.Backend = backend
	return ar
}

func (ar *AudioReceiver) BeginReception() (err error) {
	if ar.Transport == nil {
		err = errors.New("Cannot begin reception before connection to transmitter is established")
		return
	}
	log.Infoln("Initiating audio reception")

	bufsize := ar.Backend.GetBufferSize()
	buf := make(SoundBytes, bufsize)
	bufBytes := []byte(buf)
	for {
		ar.Lock()
		log.Debugf("Attempting to read data of size %d bytes", bufsize)
		if _, err = ar.Read(bufBytes); err != nil {
			ar.Unlock()
			return
		}
		if ar.ReceptionCallback != nil {
			ar.ReceptionCallback(&bufBytes)
		}
		var ret int
		if ret, err = ar.Backend.Write(bufBytes); ret < 0 {
			ar.Unlock()
			err = errors.New(fmt.Sprintf("Failed to write data to %s: %v", ar.ApiType.ApiString(), ret))
			return
		}
		ar.Unlock()
	}
	return
}

func (ar *AudioReceiver) Listen(proto string, addr string, callback func(transport Transport)) (err error) {
	log.Info("Listening for a connection...")

	cb := func(transport Transport) {
		ar.Transport = transport
		log.Info("Received connection from:", ar.Transport.RemoteAddr())
		callback(transport)
	}

	switch proto {
	case "tcp":
		server := NewTCPServer()
		if err = server.Bind(addr); err != nil {
			return
		}
		if err = server.Listen(cb); err != nil {
			return
		}
	case "kcp":
		/*
			if kcpListener, err = kcp.ListenWithOptions(addr, nil, 10, 3); err != nil {
				return
			} else {
				conn, err = kcpListener.AcceptKCP()
			}
			transport := &BaseTransport{}
			transport.Conn = conn
			ar.Transport = transport
		*/
	case "udp":
		server := NewUDPServer()
		if err = server.Listen(cb); err != nil {
			return
		}
	}
	return
}
