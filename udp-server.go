package audiotransport

import (
	"fmt"

	"github.com/xtaci/kcp-go"
)

type UdpServer struct {
	*Transport
}

func NewUDPServer() *UdpServer {
	server := UdpServer{}
	server.Transport = &Transport{}
	return &server
}

func (server *UdpServer) Listen(addr string) error {
	listener, err := kcp.ListenWithOptions(addr, nil, 3, 10)
	if err != nil {
		return err
	}
	conn, err := listener.AcceptKCP()
	if err != nil {
		return err
	}
	fmt.Println("Received connection from: %v", conn.RemoteAddr())
	server.Conn = conn
	return err
}
