package audiotransport

import (
	"encoding/binary"
	"fmt"
	"net"
	"reflect"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func write(conn net.Conn, b []byte) (n int, err error) {
	size := make([]byte, 4)
	binary.LittleEndian.PutUint32(size, uint32(len(b)))
	if _, err = conn.Write(size); err != nil {
		return
	}
	return conn.Write(b)
}

func read(conn net.Conn) (dataBytes []byte, err error) {
	sizeBytes := make([]byte, 4)
	if _, err = conn.Read(sizeBytes); err != nil {
		return
	}
	size := binary.LittleEndian.Uint32(sizeBytes)
	dataBytes = make([]byte, size)
	if _, err = conn.Read(dataBytes); err != nil {
		return
	}
	return
}

func TestServer(t *testing.T) {
	assert := assert.New(t)

	server := NewUDPServer()
	err := server.Listen("127.0.0.1", 6654)
	assert.Nil(err, "Failed to start server")

	conn, err := net.Dial("udp", "127.0.0.1:6654")
	assert.Nil(err, "Failed to connect to server")

	CommonTest(assert, server.Transport, conn)
}

func CommonTest(assert *assert.Assertions, transport *Transport, conn net.Conn) {
	go func() {
		_, err := write(conn, []byte("MAGIC-CLIENT"))
		assert.Nil(err, "Failed to write data to transport", err)

		_, err = write(conn, []byte(MAGIC))
		assert.Nil(err, "Failed to write MAGIC string to transport", err)

		err = conn.Close()
		assert.Nil(err, "Failed to close connection to transport", err)
	}()

	dataChan := make(chan []byte)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		transport.AsyncRead(dataChan)
		close(dataChan)
	}()

	data := make([]byte, 0)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for b := range dataChan {
			data = append(data, b...)
		}
	}()

	wg.Wait()
	assert.Equal(string(data), "MAGIC-CLIENT", "Data does not match")

	err := transport.Close()
	assert.Nil(err, fmt.Sprintf("Failed to close transport: %v", err))

}

func TransportData(assert *assert.Assertions, data []byte, port int) {
	var err error

	server := NewUDPServer()
	client := NewUDPClient()

	go func() {
		addr := fmt.Sprintf("127.0.0.1:%v", port)
		err = client.Connect(addr)
		assert.Nil(err, "Failed to connect to server")

		_, err = client.WriteBytes(data)
		assert.Nil(err, "Failed to write bytes to server", err)

		//_, err = client.WriteBytes([]byte(MAGIC))
		//client.Close()
	}()

	err = server.Listen("127.0.0.1", port)
	assert.Nil(err, "Failed to listen for connections", err)

	rcvData, err := server.ReadBytes()
	assert.Nil(err, "Failed to read bytes from client", err)
	assert.True(reflect.DeepEqual(data, rcvData), "Data does not match")

	//data, err = server.ReadBytes()
	//assert.Equal(MAGIC, string(data))
}

func TestTransport(t *testing.T) {
	assert := assert.New(t)
	TransportData(assert, []byte("TEST-TRANSPORT-STRING"), 6554)
}

func TestTransportLargeFile(t *testing.T) {
	assert := assert.New(t)

	_, data, err := ParseWavFile("test.wav")
	assert.Nil(err, "Failed to parse WAV file", err)

	TransportData(assert, data, 6555)
}
