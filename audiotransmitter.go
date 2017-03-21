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

type AudioTransmitter struct {
	Transport
	ApiType
	sync.Mutex
	Name                 string
	Device               string
	CaptureIdx           int32
	samplerate           int32
	channels             int32
	TransmissionCallback func(b []byte, len int32)
}

func NewAudioTransmitter(apiType ApiType, name string, device string, samplerate int32, channels int32) *AudioTransmitter {
	var idx int32

	switch apiType {
	case ALSA_API:
		idx = alsa.Alsa_init(device, samplerate, channels, 0)
	case PULSE_API:
		idx = alsa.Pa_init(name, device, samplerate, channels, 0)
	}
	if idx < 0 {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Failed to initialize %s: %v", apiType.ApiString(), idx))
		return nil
	}

	at := &AudioTransmitter{}
	at.ApiType = apiType
	at.Name = name
	at.Device = device
	at.CaptureIdx = idx
	at.initialize(samplerate, channels)
	return at
}

func (at *AudioTransmitter) initialize(samplerate int32, channels int32) {
	at.samplerate = samplerate
	at.channels = channels
}

func (at *AudioTransmitter) BeginTransmission() (err error) {
	if at.Transport == nil {
		err = errors.New("Cannot begin transmission before connection to receiver is established")
		return
	}
	bufsize := at.GetBufferSize(at.samplerate, at.channels)
	buf := make([]byte, bufsize)
	for {
		switch at.ApiType {
		case ALSA_API:
			alsa.Alsa_readi(at.CaptureIdx, &buf, bufsize)
		case PULSE_API:
			alsa.Pa_handle_read(at.CaptureIdx, &buf, bufsize)
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
