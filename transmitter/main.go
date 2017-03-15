package main

/*
#cgo LDFLAGS: -L.. -lalsa -lasound -lpulse -lpulse-simple
*/
import "C"

import (
	"fmt"
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/gurupras/audiotransport"
)

var (
	addr  *string
	proto *string
)

func setupParser() {
	addr = kingpin.Arg("receiver-address", "Address of receiver").Required().String()
	proto = kingpin.Flag("protocol", "tcp/udp").Short('P').Default("udp").String()
}
func main() {
	setupParser()
	kingpin.Parse()
	var err error

	audioTransmitter := audiotransport.NewAudioTransmitter("transmitter", "alsa_output.pci-0000_00_05.0.analog-stereo.monitor", 48000, 2)
	if err = audioTransmitter.Connect(*proto, *addr); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Failed to connet to server: %v", err))
		return
	}
	fmt.Println("Connected to remote receiver:", audioTransmitter.RemoteAddr())

	audioTransmitter.BeginTransmission()
}
