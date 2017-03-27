package audiotransport

import (
	"fmt"
	"net"
)

type TCPServer struct {
	*Server
	listener *net.TCPListener
}

type TCPSession struct {
	*net.TCPConn
	RemoteTCPAddr *net.TCPAddr
}

type TCPTransport struct {
	*BaseTransport
	*TCPSession
}

func NewTCPServer() *TCPServer {
	server := TCPServer{}
	server.Server = &Server{}
	return &server
}

func (server *TCPServer) Bind(addr string) (err error) {
	var tcpAddr *net.TCPAddr

	if tcpAddr, err = net.ResolveTCPAddr("tcp", addr); err != nil {
		return
	}

	if server.listener, err = net.ListenTCP("tcp", tcpAddr); err != nil {
		return
	}
	return
}

func (server *TCPServer) Listen(callback func(transport Transport)) (err error) {
	var conn *net.TCPConn
	go func() {
		_ = <-server.signalChan
		server.listener.Close()
	}()
	for {
		if conn, err = server.listener.AcceptTCP(); err != nil {
			return
		}

		session := &TCPSession{}
		session.TCPConn = conn
		transport := &TCPTransport{}
		transport.BaseTransport = &BaseTransport{conn}
		transport.TCPSession = session
		fmt.Printf("Received connection from: %v\n", conn.RemoteAddr())
		callback(transport)
	}
	return
}
