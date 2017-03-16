package audiotransport

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransmitter(t *testing.T) {
	assert := assert.New(t)

	go func() {
		at := NewAudioTransmitter(PULSE_API, "TestTransmitter", "alsa_output.pci-0000_00_05.0.analog-stereo.monitor", 48000, 2)
		assert.NotNil(at, "Failed to initialize audio transmitter")
		err := at.Connect("udp", "127.0.0.1:6556")
		assert.Nil(err, "Failed to connect to server", err)
		err = at.BeginTransmission()
	}()

	// Now start a dumb udp server that discards the data
	wg := sync.WaitGroup{}
	wg.Add(1)
	runOnce := false

	server := NewUDPServer()
	callback := func(transport Transport) {
		for {
			_, _ = transport.ReadBytes()
			if runOnce == false {
				// We received data..terminate
				wg.Done()
				runOnce = true
			}
		}
	}

	err := server.Listen("127.0.0.1:6556", callback)
	assert.Nil(err, "Failed to listen on server", err)
	wg.Wait()
}

func TestReceiver(t *testing.T) {
	assert := assert.New(t)

	go func() {
		// Start a UDP client and feed in a WAV file
		client := NewUDPClient()
		transport, err := client.Connect("127.0.0.1:6557")
		assert.Nil(err, "Failed to connect to receiver", err)

		_, data, err := ParseWavFile("test.wav")
		assert.Nil(err, "Failed to parse WAV file", err)

		transport.WriteBytes(data)
		assert.Nil(err, "Failed to write bytes to server", err)
		fmt.Println("Finished writing WAV data:", len(data))
	}()

	wg := sync.WaitGroup{}
	wg.Add(1)
	callback := func(b *[]byte) {
		wg.Done()
	}

	go func() {
		ar := NewAudioReceiver(PULSE_API, "TestReceiver", "alsa_output.pci-0000_00_05.0.analog-stereo", 48000, 2)
		assert.NotNil(ar, "Failed to initialize audio receiver")
		err := ar.Listen("udp", "127.0.0.1:6557")
		assert.Nil(err, "Failed to listen for connections", err)
		err = ar.BeginReception(callback)
		assert.Nil(err, "Failed to receive data from receiver", err)
	}()

	wg.Wait()
}
