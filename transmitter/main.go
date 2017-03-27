package main

import "C"

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/alecthomas/kingpin"
	"github.com/gurupras/audiotransport"
)

var (
	name          *string
	addr          *string
	proto         *string
	device        *string
	api           *string
	filterSilence *bool
	verbose       *bool
)

func setupParser() {
	name = kingpin.Flag("name", "program name. This is used as filename in FILE method").Short('n').Default("transmitter").String()
	addr = kingpin.Arg("receiver-address", "Address of receiver").Required().String()
	proto = kingpin.Flag("protocol", "tcp/udp").Short('P').Default("udp").String()
	device = kingpin.Flag("device", "Device from which to capture and transmit").Short('d').String()
	api = kingpin.Flag("method", "Which mechanism to use.. ALSA/PULSE").Short('m').Default("PULSE").String()
	filterSilence = kingpin.Flag("filter-silence", "Filter out empty audio data").Short('f').Default("true").Bool()
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
			dev = "hw:1,1"
		case audiotransport.PULSE_API:
			dev = "alsa_output.pci-0000_00_05.0.analog-stereo.monitor"
		}
	} else {
		dev = *device
	}
	log.Debugf("Device=%s\n", dev)

	addrs := strings.Split(*addr, ",")

	wg := sync.WaitGroup{}
	wg.Add(len(addrs))

	latencyGoroutine := func(addr string, audioTransmitter *audiotransport.AudioTransmitter) {
		for {
			lat, _ := audioTransmitter.Backend.GetLatency()
			log.Infof("%v: Transmitter latency=%0.0f", addr, float64(lat))
			time.Sleep(1000 * time.Millisecond)
		}
	}

	for _, addr := range addrs {
		audioTransmitter := audiotransport.NewAudioTransmitter(apiType, *name, dev, 48000, 2, true)
		if err = audioTransmitter.Connect(*proto, addr); err != nil {
			fmt.Fprintln(os.Stderr, fmt.Sprintf("Failed to connet to server: %v", err))
			return
		}
		go latencyGoroutine(addr, audioTransmitter)
		log.Infoln("Connected to remote receiver:", addr)

		var lastTime time.Time = time.Now()
		var size uint32 = 0
		audioTransmitter.TransmissionCallback = func(transport audiotransport.Transport, b []byte, len uint32) {
			now := time.Now()
			if now.Sub(lastTime).Seconds() < 1.0 {
				size += len
			} else {
				log.Infof("%v: Bandwidth: %0.2fKBps\n", transport, float32(size)/1024.0)
				size = 0
				lastTime = now
			}
		}
		go func() {
			defer wg.Done()
			if err = audioTransmitter.BeginTransmission(); err != nil {
				fmt.Fprintln(os.Stderr, fmt.Sprintf("Transmission failed with error: %v", err))
			}
		}()
	}
	wg.Wait()
}
