package audiotransport

import "github.com/xtaci/kcp-go"

type UdpClient struct {
	*Transport
}

func NewUDPClient() *UdpClient {
	client := &UdpClient{}
	client.Transport = &Transport{}
	return client
}

func (client *UdpClient) Connect(addr string) error {
	conn, err := kcp.DialWithOptions(addr, nil, 3, 10)
	if err != nil {
		return err
	}
	client.Conn = conn
	return err
}
