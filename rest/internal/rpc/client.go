package rpc

import (
	"fmt"
	"time"

	"github.com/zhlii/wechat-box/rest/internal/logs"
)

type Client struct {
	CmdClient *CmdClient
	MsgClient *MsgClient
	Usr       *UserInfo
	Contacts  []*RpcContact
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
	err := c.CmdClient.socket.conn(5)
	if err != nil {
		return err
	}

	go func(client *Client) {
		for {
			isLogin, err := client.CmdClient.IsLogin()

			if err != nil {
				logs.Error(err.Error())
				time.Sleep(10 * time.Second)
				continue
			}

			if isLogin {
				break
			}

			logs.Info("not login, sleeping for 5s...")
			time.Sleep(5 * time.Second)
		}

		usr, err := client.CmdClient.GetSelfInfo()
		if err != nil {
			logs.Error(fmt.Sprintf("get current user failed. %v", err))
		} else {
			client.Usr = usr
		}

		c.FreshContacts()
	}(c)

	return nil
}

func (s *Client) FreshContacts() {
	contacts, err := s.CmdClient.GetContacts()

	if err != nil {
		logs.Error(fmt.Sprintf("get contacts error: %v", err))
		s.Contacts = []*RpcContact{}
	} else {
		s.Contacts = contacts
		logs.Debug(fmt.Sprintf("get %d contacts", len(s.Contacts)))
	}
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
