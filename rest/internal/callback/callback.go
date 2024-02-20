package callback

import (
	"github.com/zhlii/wechat-box/rest/internal/config"
	"github.com/zhlii/wechat-box/rest/internal/rpc"
)

type Handler struct {
	Callback func(c *rpc.Client, msg *rpc.WxMsg)
}

var handlers = make(map[string]*Handler)

func Setup() []*Handler {
	handlerLog()
	handlerEcho()
	handlerAutoAcceptFriendInvite()
	handlerSaveMessage()
	handlerSpark()

	list := make([]*Handler, 0)

	for k, v := range handlers {
		cfg := config.Data.Callbacks[k]

		if cfg["enable"] == "true" {
			list = append(list, v)
		}
	}

	return list
}
