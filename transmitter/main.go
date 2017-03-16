package main

/*
#cgo LDFLAGS: -L.. -lalsa -lasound -lpulse -lpulse-simple
*/
import "C"

import (
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/kingpin"
	"github.com/gurupras/audiotransport"
)

var (
	addr   *string
	proto  *string
	device *string
	api    *string
)

func setupParser() {
	addr = kingpin.Arg("receiver-address", "Address of receiver").Required().String()
	proto = kingpin.Flag("protocol", "tcp/udp").Short('P').Default("udp").String()
	device = kingpin.Flag("device", "Device from which to capture and transmit").Short('d').String()
	api = kingpin.Flag("method", "Which mechanism to use.. ALSA/PULSE").Short('m').Default("PULSE").String()
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

	audioTransmitter := audiotransport.NewAudioTransmitter(apiType, "transmitter", dev, 48000, 2)
	if err = audioTransmitter.Connect(*proto, *addr); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Failed to connet to server: %v", err))
		return
	}
	fmt.Println("Connected to remote receiver:", audioTransmitter.RemoteAddr())

	audioTransmitter.BeginTransmission()
}
