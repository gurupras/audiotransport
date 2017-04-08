package main

import (
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/alecthomas/kingpin"
	"github.com/gurupras/audiotransport"
	"github.com/gurupras/audiotransport/pacmd"
)

func main() {
	app := kingpin.New("receiver", "Receive audio")
	config := audiotransport.ParseArgs(app, os.Args)

	if strings.Compare(config.Device, "") == 0 {
		log.Infof("No device specified... Using defaults")
		switch config.Api {
		case audiotransport.ALSA_API:
			config.Device = "hw:1,1"
		case audiotransport.PULSE_API:
			devices, err := pacmd.ListSinks()
			if err != nil {
				log.Fatalf("Failed to list pacmd sources: %v", err)
			}
			config.Device = devices[0]
		}
	}
	log.Infof("Device=%s\n", config.Device)

	audioReceiver := audiotransport.NewAudioReceiver(config.Api, config.Name, config.Device, config.Samplerate, config.Channels)
	if audioReceiver == nil {
		panic("Failed to start receiver")
	}

	audioReceiver.ReceptionCallback = func(data *[]byte) (err error) {
		log.Debugf("Wrote %d bytes", len(*data))
		return
	}

	callback := func(transport audiotransport.Transport) {
		if err := audioReceiver.BeginReception(); err != nil {
			log.Warnln(err)
		}
	}

	for _, addr := range config.Addrs {
		if err := audioReceiver.Listen(config.Proto, addr, callback); err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}
	}
}
