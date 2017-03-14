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

func (server *UdpServer) Listen(host string, port int) error {
	addr := fmt.Sprintf("%s:%d", host, port)
	listener, err := kcp.ListenWithOptions(addr, nil, 3, 10)
	if err != nil {
		return err
	}
	conn, err := listener.Accept()
	if err != nil {
		return err
	}
	server.Conn = conn
	return err
}
