package audiotransport

import (
	"net"

	"github.com/xtaci/kcp-go"
)

type UdpClient struct {
}

func NewUDPClient() *UdpClient {
	client := &UdpClient{}
	return client
}

func (client *UdpClient) Connect(addr string) (transport Transport, err error) {
	var udpAddr *net.UDPAddr
	var conn *net.UDPConn

	if false {
		_, err = kcp.DialWithOptions(addr, nil, 3, 10)
		if err != nil {
			return
		}
		_ = conn
	}

	if udpAddr, err = net.ResolveUDPAddr("udp", addr); err != nil {
		return
	}
	if conn, err = net.DialUDP("udp", nil, udpAddr); err != nil {
		return
	}

	session := &UDPSession{}
	session.UDPConn = conn
	session.RemoteUDPAddr = udpAddr

	udpTransport := &UDPTransport{}
	udpTransport.UDPSession = session

	transport = udpTransport

	return
}
