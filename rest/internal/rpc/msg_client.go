package rpc

import (
	"fmt"

	"github.com/zhlii/wechat-box/rest/internal/logs"
)

type MsgCallback func(msg *WxMsg)

type MsgClient struct {
	socket    *protobufferSocket
	callbacks map[string]MsgCallback
}

func (c *MsgClient) Close(ks ...string) error {
	logs.Debug("close msg client")
	if len(c.callbacks) > 0 && len(ks) > 0 {
		for _, k := range ks {
			delete(c.callbacks, k)
		}
		if len(c.callbacks) > 0 {
			return nil
		}
	}
	// 关闭消息推送
	c.callbacks = nil
	return c.socket.close()
}

// 创建消息接收器
// param cb MsgCallback 消息回调函数
// return string 接收器唯一标识
func (c *MsgClient) Register(cb MsgCallback) (string, error) {
	key := Rand(16)
	if c.callbacks == nil {
		if err := c.socket.conn(0); err != nil {
			logs.Error("msg socket conn error")
			return "", err
		}
		c.callbacks = map[string]MsgCallback{
			key: cb,
		}
		go func() {
			for len(c.callbacks) > 0 {
				if resp, err := c.socket.recv(); err == nil {
					msg := resp.GetWxmsg()
					for _, f := range c.callbacks {
						go f(msg)
					}
				} else {
					logs.Warn(fmt.Sprintf("msg receiver error: %v", err))
				}
			}
		}()
	} else {
		c.callbacks[key] = cb
	}
	return key, nil
}
