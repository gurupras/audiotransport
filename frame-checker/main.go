package main

/*
#cgo LDFLAGS: -L.. -lalsa -lasound -lpulse -lpulse-simple
*/
import "C"

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/alecthomas/kingpin"
	"github.com/gurupras/audiotransport"
)

var (
	addr      *string
	proto     *string
	device    *string
	api       *string
	frameSize *int
)

type Frame struct {
	Bytes []byte
	Size  int
}

func NewFrame(size int) *Frame {
	f := &Frame{}
	f.Bytes = make([]byte, size)
	f.Size = size
	return f
}

type FrameArray []*Frame

func (fa FrameArray) String() string {
	b := bytes.NewBuffer(nil)
	for i := 0; i < len(fa); i++ {
		b.WriteString(fa[i].String() + " ")
	}
	return b.String()
}

func (f *Frame) FromReader(reader io.Reader) {
	reader.Read(f.Bytes)
}

func (f *Frame) String() string {
	b := bytes.NewBuffer(nil)
	b.WriteString("[")
	for i := 0; i < f.Size; i++ {
		b.WriteString(fmt.Sprintf("%03d", f.Bytes[i]))
		if i < f.Size-1 {
			b.WriteString(" ")
		}
	}
	b.WriteString("]")
	return b.String()
}

func ParseFrame(frameSize int, reader io.Reader) *Frame {
	f := NewFrame(frameSize)
	f.FromReader(reader)
	return f
}

func setupParser() {
	proto = kingpin.Flag("protocol", "tcp/udp").Short('P').Default("udp").String()
	device = kingpin.Flag("device", "Device from which to capture and transmit").Short('d').String()
	api = kingpin.Flag("method", "Which mechanism to use.. ALSA/PULSE").Short('m').Default("PULSE").String()
	frameSize = kingpin.Flag("frame-size", "Size of each frame").Short('s').Default("2").Int()
}

func main() {
	setupParser()
	kingpin.Parse()
	var err error

	var apiType audiotransport.ApiType
	switch *api {
	case "PULSE":
		apiType = audiotransport.PULSE_API
	case "ALSA":
		apiType = audiotransport.ALSA_API
	default:
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Invalid API: %v", *api))
		return
	}

	var dev string
	if device == nil || strings.Compare(*device, "") == 0 {
		switch apiType {
		case audiotransport.ALSA_API:
			dev = "hw:1,1"
		case audiotransport.PULSE_API:
			dev = "alsa_output.pci-0000_00_05.0.analog-stereo.monitor"
		}
	} else {
		dev = *device
	}
	log.Debugf("Device=%s\n", dev)

	audioTransmitter := audiotransport.NewAudioTransmitter(apiType, "transmitter", dev, 48000, 2)

	if err = audioTransmitter.Connect(*proto, "255.255.255.255:6554"); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Failed to connet to server: %v", err))
		return
	}
	log.Infoln("Connected to remote receiver:", audioTransmitter.RemoteAddr())

	audioTransmitter.TransmissionCallback = func(b []byte, len int32) {
		sum := int64(0)

		for i := int32(0); i < len; i++ {
			sum += int64(b[i])
		}
		if sum > 0 {
			// There was some data
			reader := bytes.NewReader(b)
			frames := FrameArray{}
			for i := int32(0); i < len/2; i++ {
				f := ParseFrame(*frameSize, reader)
				frames = append(frames, f)
			}
			fmt.Println(frames.String())
		}
	}
	audioTransmitter.BeginTransmission()
}
