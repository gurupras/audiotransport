package audiotransport

import (
	"encoding/binary"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransmitter(t *testing.T) {
	assert := assert.New(t)

	at := NewAudioTransmitter("TestTransmitter", "alsa_output.pci-0000_00_05.0.analog-stereo.monitor", 48000, 2)
	assert.NotNil(at, "Failed to initialize audio transmitter")

	go func() {
		err := at.Connect("127.0.0.1:6556")
		assert.Nil(err, "Failed to connect to server", err)
		err = at.BeginTransmission()
	}()

	// Now start a dumb udp server that discards the data
	server := NewUDPServer()
	err := server.Listen("127.0.0.1:6556")
	assert.Nil(err, "Failed to listen on server", err)
	for {
		_, _ = server.ReadBytes()
		// We received data..terminate
		return
	}
}

func TestReceiver(t *testing.T) {
	assert := assert.New(t)

	go func() {
		// Start a UDP client and feed in a WAV file
		client := NewUDPClient()
		err := client.Connect("127.0.0.1:6557")
		assert.Nil(err, "Failed to connect to receiver", err)

		_, data, err := ParseWavFile("test.wav")
		assert.Nil(err, "Failed to parse WAV file", err)

		lengthBytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(lengthBytes, uint32(len(data)))
		client.Write(lengthBytes)
		client.Write(data)
		assert.Nil(err, "Failed to write bytes to server", err)
		fmt.Println("Finished writing WAV data:", len(data))
	}()

	wg := sync.WaitGroup{}
	wg.Add(1)
	callback := func(b *[]byte) {
		wg.Done()
	}

	go func() {
		ar := NewAudioReceiver("TestReceiver", "alsa_output.pci-0000_00_05.0.analog-stereo", 48000, 2)
		assert.NotNil(ar, "Failed to initialize audio receiver")
		err := ar.Listen("127.0.0.1:6557")
		assert.Nil(err, "Failed to listen for connections", err)
		err = ar.BeginReception(callback)
		assert.Nil(err, "Failed to receive data from receiver", err)
	}()

	wg.Wait()
}
