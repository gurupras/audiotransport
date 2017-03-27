package audiotransport

import (
	"net"
)

type TCPClient struct {
}

func NewTCPClient() *TCPClient {
	client := &TCPClient{}
	return client
}

func (client *TCPClient) Connect(addr string) (transport Transport, err error) {
	var tcpAddr *net.TCPAddr
	var conn *net.TCPConn

	if tcpAddr, err = net.ResolveTCPAddr("tcp", addr); err != nil {
		return
	}

	if conn, err = net.DialTCP("tcp", nil, tcpAddr); err != nil {
		return
	}

	session := &TCPSession{}
	session.TCPConn = conn
	session.RemoteTCPAddr = tcpAddr

	tcpTransport := &TCPTransport{}
	tcpTransport.BaseTransport = &BaseTransport{conn}
	tcpTransport.TCPSession = session

	transport = tcpTransport
	return
}
