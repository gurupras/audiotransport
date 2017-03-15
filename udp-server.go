package audiotransport

import (
	"fmt"
	"net"

	"github.com/xtaci/kcp-go"
)

type UDPServer struct {
}

type UDPSession struct {
	*net.UDPConn
	RemoteUDPAddr *net.UDPAddr
}

type UDPTransport struct {
	*UDPSession
}

func NewUDPServer() *UDPServer {
	server := UDPServer{}
	return &server
}

func (server *UDPServer) Listen(addr string, callback func(transport Transport)) (err error) {
	if false {
		listener, err := kcp.ListenWithOptions(addr, nil, 3, 10)
		if err != nil {
			return err
		}
		conn, err := listener.AcceptKCP()
		if err != nil {
			return err
		}
		_ = conn
	}
	var udpAddr *net.UDPAddr
	var conn *net.UDPConn

	if udpAddr, err = net.ResolveUDPAddr("udp", addr); err != nil {
		return
	}

	if conn, err = net.ListenUDP("udp", udpAddr); err != nil {
		return
	}

	session := &UDPSession{}
	session.UDPConn = conn
	transport := &UDPTransport{}
	transport.UDPSession = session
	//if transport.RemoteUDPAddr, err = net.ResolveUDPAddr("udp", conn.RemoteAddr().String()); err != nil {
	//}
	fmt.Printf("Received connection from: %v\n", conn.RemoteAddr())
	callback(transport)
	return err
}
