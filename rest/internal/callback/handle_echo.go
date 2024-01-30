package callback

import "github.com/zhlii/wechat-box/rest/internal/rpc"

func handlerEcho() {
	handlers["echo"] = &Handler{
		Callback: func(c *rpc.Client, msg *rpc.WxMsg) {
			switch msg.Type {
			case 1:
				if rpc.ContactType(msg.Sender) == "文件传输助手" {
					c.CmdClient.SendTxt(msg.Content, msg.Sender, "")
				}
			}
		},
	}
}