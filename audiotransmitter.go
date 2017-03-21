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

type AudioTransmitter struct {
	Transport
	ApiType
	Backend BackendInterface
	sync.Mutex
	TransmissionCallback func(b []byte, len int32)
	FilterSilence        bool
}

type SoundBytes []byte

func (sb SoundBytes) HasData() bool {
	sum := uint64(0)
	for idx := 0; idx < len(sb); idx++ {
		sum += uint64(sb[idx])
	}
	return sum > 0
}

func NewAudioTransmitter(apiType ApiType, name string, device string, samplerate int32, channels int32, filterSilence bool) *AudioTransmitter {
	var backend BackendInterface

	switch apiType {
	case ALSA_API:
		backend = &AlsaBackend{}
	case PULSE_API:
		backend = &PulseBackend{}
	case FILE_API:
		backend = &FileBackend{}
	}
	if err := backend.Init(name, device, samplerate, channels, 0); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Failed to initialize %s: %v", apiType.ApiString(), err))
		return nil
	}
	at := &AudioTransmitter{}
	at.ApiType = apiType
	at.Backend = backend
	at.FilterSilence = filterSilence
	return at
}

func (at *AudioTransmitter) BeginTransmission() (err error) {
	if at.Transport == nil {
		err = errors.New("Cannot begin transmission before connection to receiver is established")
		return
	}
	bufsize := at.Backend.GetBufferSize()
	buf := make(SoundBytes, bufsize)
	bufBytes := []byte(buf)
	for {
		at.Backend.Read(bufBytes, bufsize)
		if at.FilterSilence && !buf.HasData() {
			// Nothing to do
			continue
		}

		at.Lock()
		log.Debugf("Attempting to send %d bytes\n", len(buf))
		if _, err = at.Write(buf); err != nil {
			err = errors.New(fmt.Sprintf("Failed to send data over transport: %v", err))
			return
		}
		if at.TransmissionCallback != nil {
			at.TransmissionCallback(buf, bufsize)
		}
		at.Unlock()
		log.Debugf("Sent %d bytes\n", len(buf))
	}
	return
}

func (at *AudioTransmitter) Connect(proto string, addr string) (err error) {
	var conn net.Conn

	at.Lock()
	defer at.Unlock()
	switch proto {
	case "tcp":
		conn, err = net.Dial("tcp", addr)
		transport := &BaseTransport{}
		transport.Conn = conn
		at.Transport = transport
	case "kcp":
		conn, err = kcp.DialWithOptions(addr, nil, 10, 3)
		transport := &BaseTransport{}
		transport.Conn = conn
		at.Transport = transport
	case "udp":
		var transport Transport
		client := NewUDPClient()
		if transport, err = client.Connect(addr); err != nil {
			return
		}
		at.Transport = transport
	}
	return
}
