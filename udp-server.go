package audiotransport

import (
	"fmt"
	"net"
)

type UDPServer struct {
	*Server
	addr *net.UDPAddr
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
	server.Server = &Server{}
	return &server
}

func (server *UDPServer) Bind(addr string) (err error) {
	server.addr, err = net.ResolveUDPAddr("udp", addr)
	return
}

func (server *UDPServer) Listen(callback func(transport Transport)) (err error) {
	/*
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
	*/
	var conn *net.UDPConn

	if conn, err = net.ListenUDP("udp", server.addr); err != nil {
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

	// We want listen to be a blocking operation since UDP is connection-less
	_ = <-server.signalChan
	return err
}
