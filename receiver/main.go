package main

/*
#cgo LDFLAGS: -L.. -lalsa -lasound -lpulse -lpulse-simple
*/
import "C"

import (
	"fmt"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/alecthomas/kingpin"
	"github.com/gurupras/audiotransport"
)

var (
	name   *string
	addr   *string
	proto  *string
	device *string
	api    *string
)

func setupParser() {
	name = kingpin.Flag("name", "program name. This is used as filename in FILE method").Short('n').Default("receiver").String()
	addr = kingpin.Arg("receiver-address", "Address of receiver").Required().String()
	proto = kingpin.Flag("protocol", "tcp/udp").Short('P').Default("udp").String()
	device = kingpin.Flag("device", "Device to use for playback").Short('d').String()
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
	case "FILE":
		apiType = audiotransport.FILE_API
	default:
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Invalid API: %v", *api))
		return
	}

	var dev string
	if device == nil || strings.Compare(*device, "") == 0 {
		switch apiType {
		case audiotransport.ALSA_API:
			dev = "default"
		case audiotransport.PULSE_API:
			dev = "NULL"
		}
	} else {
		dev = *device
	}
	log.Debugf("Device=%s\n", dev)

	audioReceiver := audiotransport.NewAudioReceiver(apiType, *name, dev, 48000, 2)

	if err = audioReceiver.Listen(*proto, *addr); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Failed to connet to server: %v", err))
		return
	}

	if err = audioReceiver.BeginReception(nil); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
