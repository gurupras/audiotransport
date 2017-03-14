package audiotransport

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"

	"github.com/xtaci/kcp-go"
)

type Transport struct {
	*kcp.UDPSession
}

const MAGIC string = "@@$!@@#@@@"

func (transport *Transport) ReadBytes() (data []byte, err error) {
	var n int
	// First, get the size of the next data frame
	sizeBytes := make([]byte, 4)
	if n, err = transport.Read(sizeBytes); err != nil {
		return
	} else if n != 4 {
		err = errors.New(fmt.Sprintf("Expected to read %v bytes. Read %v", 4, n))
		return
	}
	size := binary.LittleEndian.Uint32(sizeBytes)
	fmt.Printf("Attempting to read %d bytes\n", size)
	data = make([]byte, size)

	n, err = transport.UDPSession.Read(data)
	fmt.Println("Read:", len(data))
	return
}

func (transport *Transport) AsyncRead(byteArrayChan chan []byte) {
	magicBytes := []byte(MAGIC)
	for {
		data, err := transport.ReadBytes()
		if err != nil {
			fmt.Fprintln(os.Stderr, fmt.Sprintf("Encountered error in AsyncRead: %v", err))
			break
		}
		if len(data) == len(magicBytes) && bytes.Equal(magicBytes, data) {
			// This is our termination signal. Don't send this.
			// Instead, just break from the loop
			break
		}
		byteArrayChan <- data
	}
}

func (transport *Transport) WriteBytes(bytes []byte) (n int, err error) {
	// The protocol is to first write the length and then the bytes itself
	lengthBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(lengthBytes, uint32(len(bytes)))
	if n, err = transport.UDPSession.Write(lengthBytes); err != nil {
		err = errors.New(fmt.Sprintf("Failed to write length on transport: %v", err))
		return
	}
	n, err = transport.UDPSession.Write(bytes)
	fmt.Printf("Wrote %d bytes\n", n)
	return
}
