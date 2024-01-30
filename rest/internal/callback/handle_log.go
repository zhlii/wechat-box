package callback

import (
	"fmt"

	"github.com/zhlii/wechat-box/rest/internal/logs"
	"github.com/zhlii/wechat-box/rest/internal/rpc"
)

func handlerLog() {
	handlers["log"] = &Handler{
		Callback: func(c *rpc.Client, msg *rpc.WxMsg) {
			logs.Debug(fmt.Sprintf("receive msg: %v", msg))
		},
	}
}
