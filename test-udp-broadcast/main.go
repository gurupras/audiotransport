package main

/*
#cgo LDFLAGS: -L.. -lalsa -lasound -lpulse -lpulse-simple
*/
import "C"

import (
	"bufio"
	"fmt"
	"os"

	"github.com/alecthomas/kingpin"
	"github.com/gurupras/audiotransport"
)

var (
	app        *kingpin.Application
	addr       *string
	clientAddr *string
	serverAddr *string
	serverCmd  *kingpin.CmdClause
	clientCmd  *kingpin.CmdClause
)

func setupParser() {
	app = kingpin.New("udp-chat", "Test application for UDP")
	serverCmd = app.Command("server", "server-mode")
	clientCmd = app.Command("client", "client-mode")
	clientAddr = clientCmd.Arg("server-address", "server address").Required().String()
}

func server(addr string) {
	server := audiotransport.NewUDPServer()
	callback := func(transport audiotransport.Transport) {
		dataChannel := make(chan []byte)
		go transport.AsyncRead(dataChannel)
		for data := range dataChannel {
			fmt.Printf("Client: %s\n", string(data))
		}
	}
	if err := server.Listen("0.0.0.0:6554", callback); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

}

func client(addr string) {
	client := audiotransport.NewUDPClient()
	transport, err := client.Connect(addr)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Failed to connect: %v", err))
		return
	}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		if _, err := transport.WriteBytes([]byte(input)); err != nil {
			fmt.Fprintln(os.Stderr, fmt.Sprintf("Failed to write to UDP socket: %v", err))
			break
		}
	}
}

func main() {
	setupParser()
	cmd, err := app.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Failed to parse arguments: %v", err))
		return
	}
	switch kingpin.MustParse(cmd, err) {
	case serverCmd.FullCommand():
		server("")
	case clientCmd.FullCommand():
		client(*clientAddr)
	}

}
