package main

import (
	"os"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/alecthomas/kingpin"
	"github.com/gurupras/audiotransport"
	"github.com/gurupras/audiotransport/pacmd"
)

func main() {
	app := kingpin.New("transmitter", "Transmit audio")
	config := audiotransport.ParseArgs(app, os.Args)

	if strings.Compare(config.Device, "") == 0 {
		log.Infof("No device specified... Using defaults")
		switch config.Api {
		case audiotransport.ALSA_API:
			config.Device = "hw:1,1"
		case audiotransport.PULSE_API:
			devices, err := pacmd.ListSources()
			if err != nil {
				log.Fatalf("Failed to list pacmd sources: %v", err)
			}
			config.Device = devices[0]
		}
	}
	log.Infof("Device=%s\n", config.Device)

	wg := sync.WaitGroup{}
	wg.Add(len(config.Addrs))

	latencyGoroutine := func(addr string, audioTransmitter *audiotransport.AudioTransmitter) {
		for {
			lat, _ := audioTransmitter.Backend.GetLatency()
			log.Infof("%v: Transmitter latency=%0.0f", addr, float64(lat))
			time.Sleep(1000 * time.Millisecond)
		}
	}

	for _, addr := range config.Addrs {
		audioTransmitter := audiotransport.NewAudioTransmitter(config.Api, config.Name, config.Device, config.Samplerate, config.Channels, true)
		if err := audioTransmitter.Connect(config.Proto, addr); err != nil {
			log.Fatalf("Failed to connet to server: %v", err)
		}
		go latencyGoroutine(addr, audioTransmitter)
		log.Infof("Connected to remote receiver: %v", addr)

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
			if err := audioTransmitter.BeginTransmission(); err != nil {
				log.Fatalf("Transmission failed with error: %v", err)
			}
		}()
	}
	wg.Wait()
}
