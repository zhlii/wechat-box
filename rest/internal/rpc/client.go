package rpc

import (
	"errors"
	"time"
)

type Client struct {
	CmdClient *CmdClient
	MsgClient *MsgClient
}

func NewClient(host string, port int) *Client {
	client := &Client{
		CmdClient: &CmdClient{
			socket: newProtobufferSocker(host, port),
		},
		MsgClient: &MsgClient{
			socket: newProtobufferSocker(host, port+1),
		},
	}

	return client
}

func (c *Client) Connect() error {
	return c.CmdClient.socket.conn(25)
}

func (c *Client) RegisterCallback(callback MsgCallback) error {
	if c.MsgClient.callbacks == nil {
		if c.CmdClient.EnableMsgReciver(true) != 0 {
			return errors.New("failed to enable msg server")
		}
	}

	time.Sleep(time.Second)

	_, err := c.MsgClient.Register(callback)

	return err
}

func (c *Client) Close() {
	c.MsgClient.close()

	c.CmdClient.DisableMsgReciver()
	c.CmdClient.Close()
}
