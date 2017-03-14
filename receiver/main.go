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
	addr *string
)

func setupParser() {
	addr = kingpin.Arg("receiver-address", "Address of receiver").Required().String()
}
func main() {
	kingpin.Parse()
	var err error

	audioReceiver := audiotransport.NewAudioReceiver("transmitter", "NULL", 48000, 2)
	if err = audioReceiver.Listen(*addr); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Failed to connet to server: %v", err))
		return
	}

	audioReceiver.BeginReception(nil)
}
