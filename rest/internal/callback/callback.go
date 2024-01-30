package callback

import "github.com/zhlii/wechat-box/rest/internal/rpc"

type Handler struct {
	Callback func(c *rpc.Client, msg *rpc.WxMsg)
}

var handlers = make(map[string]*Handler)

func Setup() []*Handler {
	handlerEcho()

	list := make([]*Handler, 0)

	for _, v := range handlers {
		list = append(list, v)
	}

	return list
}
