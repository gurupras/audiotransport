package audiotransport

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"os"
)

type Transport interface {
	net.Conn
	ReadBytes() ([]byte, error)
	AsyncRead(chan []byte)
	WriteBytes([]byte) (int, error)
	String() string
}

type BaseTransport struct {
	net.Conn
}

const MAGIC string = "@@$!@@#@@@"

func (transport *BaseTransport) ReadBytes() (data []byte, err error) {
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

	n, err = transport.Conn.Read(data)
	fmt.Println("Read:", len(data))
	return
}

func (transport *BaseTransport) AsyncRead(byteArrayChan chan []byte) {
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

func (transport *BaseTransport) WriteBytes(data []byte) (n int, err error) {
	// The protocol is to first write the length and then the bytes itself
	lengthBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(lengthBytes, uint32(len(data)))
	if n, err = transport.Conn.Write(lengthBytes); err != nil {
		err = errors.New(fmt.Sprintf("Failed to write length on transport: %v", err))
		return
	}
	n, err = transport.Conn.Write(data)
	fmt.Printf("Wrote %d bytes\n", n)
	return
}

func (transport *BaseTransport) String() string {
	return fmt.Sprintf("%v", transport.RemoteAddr())
}

func (transport *UDPTransport) ReadBytes() (data []byte, err error) {
	var n int
	var remoteAddr *net.UDPAddr
	// First, get the size of the next data frame
	sizeBytes := make([]byte, 4)
	if n, remoteAddr, err = transport.ReadFromUDP(sizeBytes); err != nil {
		return
	} else if n != 4 {
		err = errors.New(fmt.Sprintf("Expected to read %v bytes. Read %v", 4, n))
		return
	}
	size := binary.LittleEndian.Uint32(sizeBytes)
	fmt.Printf("Attempting to read %d bytes\n", size)
	data = make([]byte, size)

	if n, remoteAddr, err = transport.ReadFromUDP(data); err != nil {
		err = errors.New(fmt.Sprintf("Failed to read data: %v", err))
		return
	}
	_ = remoteAddr
	fmt.Println("Read:", len(data))
	return
}

func (transport *UDPTransport) AsyncRead(byteArrayChan chan []byte) {
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

func (transport *UDPTransport) WriteBytes(data []byte) (n int, err error) {
	// The protocol is to first write the length and then the bytes itself
	lengthBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(lengthBytes, uint32(len(data)))
	if n, err = transport.Write(lengthBytes); err != nil {
		err = errors.New(fmt.Sprintf("Failed to write length on transport: %v", err))
		return
	}
	if n, err = transport.Write(data); err != nil {
		err = errors.New(fmt.Sprintf("Failed to write length on transport: %v", err))
		return
	}
	fmt.Printf("Wrote %d bytes\n", n)
	return
}

func (transport *UDPTransport) String() string {
	return fmt.Sprintf("%v", transport.UDPSession.UDPConn.RemoteAddr())
}
