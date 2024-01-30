package callback

import (
	"fmt"
	"time"

	"github.com/zhlii/wechat-box/rest/internal/logs"
	"github.com/zhlii/wechat-box/rest/internal/rpc"
)

func handlerDownloadFile() {
	handlers["download_file"] = &Handler{
		Callback: func(c *rpc.Client, msg *rpc.WxMsg) {
			switch msg.Type {
			case 3:
				if status, err := c.CmdClient.DownloadAttach(msg.Id, "", msg.Extra); status != 0 || err != nil {
					logs.Error(fmt.Sprintf("failed to download attach. status:%d error:%v", status, err))
				}

				cnt := 0
				for cnt <= 10 {
					logs.Debug(fmt.Sprintf("retry %d times", cnt))
					if path, err := c.CmdClient.DecryptImage(msg.Extra, "c:/users/root/Documents"); err != nil {
						logs.Error(fmt.Sprintf("failed to decrypt image. path:%s error:%v", path, err))
						break
					} else {
						if path != "" {
							logs.Debug(path)
							break
						} else {
							time.Sleep(time.Second)
							cnt++
						}
					}
				}
			}
		},
	}
}
