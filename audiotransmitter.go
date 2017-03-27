package audiotransport

import (
	"errors"
	"fmt"
	"net"
	"os"
	"sync"

	log "github.com/Sirupsen/logrus"
)

type AudioTransmitter struct {
	Transports []Transport
	ApiType
	Backend BackendInterface
	sync.Mutex
	TransmissionCallback func(transport Transport, b []byte, len uint32)
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

func NewAudioTransmitter(apiType ApiType, name string, device string, samplerate, channels uint32, filterSilence bool) *AudioTransmitter {
	var backend BackendInterface

	switch apiType {
	case ALSA_API:
		backend = &AlsaBackend{}
	case PULSE_API:
		backend = &PulseBackend{}
	case FILE_API:
		backend = &FileBackend{}
	}
	if err := backend.Init(name, device, samplerate, channels, false); err != nil {
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
	if len(at.Transports) == 0 {
		err = errors.New("Cannot begin transmission before connection to receiver is established")
		return
	}
	bufsize := at.Backend.GetBufferSize()
	buf := make(SoundBytes, bufsize)
	bufBytes := []byte(buf)
	for {
		at.Backend.Read(bufBytes)
		if at.FilterSilence && !buf.HasData() {
			// Nothing to do
			continue
		}

		at.Lock()
		for _, transport := range at.Transports {
			log.Debugf("%v: Attempting to send %d bytes\n", transport, len(buf))
			if _, err = transport.Write(buf); err != nil {
				err = errors.New(fmt.Sprintf("Failed to send data over transport: %v", err))
				at.Unlock()
				return
			}
			if at.TransmissionCallback != nil {
				at.TransmissionCallback(transport, buf, bufsize)
			}
			log.Debugf("Sent %d bytes\n", len(buf))
		}
		at.Unlock()
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
		at.Transports = append(at.Transports, transport)
	case "kcp":
	case "udp":
		var transport Transport
		client := NewUDPClient()
		if transport, err = client.Connect(addr); err != nil {
			return
		}
		at.Transports = append(at.Transports, transport)
	}
	return
}
