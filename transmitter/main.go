package main

/*
#cgo LDFLAGS: -L.. -lalsa -lasound -lpulse -lpulse-simple
*/
import "C"

import (
	"fmt"
	"os"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/alecthomas/kingpin"
	"github.com/gurupras/audiotransport"
	"github.com/gurupras/audiotransport/alsa"
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
	log.Debugf("Device=%s\n", dev)

	audioTransmitter := audiotransport.NewAudioTransmitter(apiType, "transmitter", dev, 96000, 2)
	go func() {
		for {
			log.Infof("Transmitter latency=%0.0f", float32(alsa.Pa_get_latency(audioTransmitter.CaptureIdx)))
			time.Sleep(1000 * time.Millisecond)
		}
	}()

	if err = audioTransmitter.Connect(*proto, *addr); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Failed to connet to server: %v", err))
		return
	}
	log.Infoln("Connected to remote receiver:", audioTransmitter.RemoteAddr())

	var lastTime time.Time = time.Now()
	var size int32 = 0
	audioTransmitter.TransmissionCallback = func(b []byte, len int32) {
		now := time.Now()
		if now.Sub(lastTime).Seconds() < 1.0 {
			size += len
		} else {
			log.Infof("Bandwidth: %0.2fKBps\n", float32(size)/1024.0)
			size = 0
			lastTime = now
		}
	}
	audioTransmitter.BeginTransmission()
}
