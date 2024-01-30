package rpc

import (
	"fmt"
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
	return c.CmdClient.socket.conn(5)
}

func (c *Client) RegisterCallback(callback MsgCallback) error {
	if c.MsgClient.callbacks == nil {
		if _, err := c.CmdClient.EnableMsgReciver(true); err != nil {
			return fmt.Errorf("failed to enable msg server. error: %v", err)
		}

		time.Sleep(time.Second)
	}

	_, err := c.MsgClient.Register(callback)

	return err
}

func (c *Client) Close() {
	c.MsgClient.Close()

	c.CmdClient.DisableMsgReciver()
	c.CmdClient.Close()
}
