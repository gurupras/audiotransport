package main

import (
	"fmt"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/alecthomas/kingpin"
	"github.com/gurupras/audiotransport"
)

var (
	name    *string
	addr    *string
	proto   *string
	device  *string
	api     *string
	verbose *bool
)

func setupParser() {
	name = kingpin.Flag("name", "program name. This is used as filename in FILE method").Short('n').Default("receiver").String()
	addr = kingpin.Arg("receiver-address", "Address of receiver").Required().String()
	proto = kingpin.Flag("protocol", "tcp/udp").Short('P').Default("udp").String()
	device = kingpin.Flag("device", "Device to use for playback").Short('d').String()
	api = kingpin.Flag("method", "Which mechanism to use.. ALSA/PULSE").Short('m').Default("PULSE").String()
	verbose = kingpin.Flag("verbose", "Enable verbose logging").Short('v').Default("false").Bool()
}
func main() {
	setupParser()
	kingpin.Parse()
	var err error

	if *verbose {
		log.SetLevel(log.DebugLevel)
		log.Debugln("Enabling verbose logging")
	}

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

	audioReceiver.ReceptionCallback = func(data *[]byte) (err error) {
		log.Debugf("Wrote %d bytes", len(*data))
		return
	}

	callback := func(transport audiotransport.Transport) {
		if err = audioReceiver.BeginReception(); err != nil {
			log.Fatalln(err)
			os.Exit(-1)
		}
	}

	if err = audioReceiver.Listen(*proto, *addr, callback); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Failed to connet to server: %v", err))
		return
	}

}
