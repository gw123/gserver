package gserver

import (
	"github.com/gw123/gserver/contracts"
	"net"
)

type Client struct {
	addr    string
	timeout int
	conn    net.Conn
	packer  contracts.DataPacker
}

func NewClient(addr string, timeout int, packer contracts.DataPacker) *Client {
	return &Client{
		addr:    addr,
		timeout: timeout,
		packer:  packer,
	}
}

func (client *Client) Connect() error {
	conn, err := net.Dial("tcp", client.addr)
	if err != nil {
		return err
	}
	// conn.SetDeadline(time.Now().Add(time.Second * time.Duration(client.timeout)))
	client.conn = conn
	return nil
}

func (client *Client) Close() error {
	return client.conn.Close()
}

func (client *Client) Send(msg *contracts.Msg) error {
	buf, err := client.packer.Pack(msg)
	if err != nil {
		return err
	}

	err = client.write(buf)
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) write(buf []byte) error {
	n := 0
	for n < len(buf) {
		num, err := client.conn.Write(buf[n:])
		if err != nil {
			return err
		}
		n += num
	}
	return nil
}

func (client *Client) Read() (*contracts.Msg, error) {
	msg, err := client.packer.UnPackFromConn(client.conn)
	if err != nil {
		return msg, err
	}
	return msg, nil
}
