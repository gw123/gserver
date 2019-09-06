package connpool

import (
	"context"
	"github.com/gw123/glog"
	"net"
)

type Client struct {
	ctx  context.Context
	conn net.Conn
}

func NewClient(ctx context.Context, conn net.Conn) *Client {
	c := &Client{
		ctx:  ctx,
		conn: conn,
	}
	return c
}

func (c *Client) GetConn() net.Conn {
	return c.conn
}

func (c *Client) Stop() {
	c.conn.Close()
}

func (c *Client) HandleRequest() error {
	glog.Debug("handle request")
	return nil
}
