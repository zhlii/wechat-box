package callback

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/zhlii/wechat-box/rest/internal/db"
	"github.com/zhlii/wechat-box/rest/internal/db/tables"
	"github.com/zhlii/wechat-box/rest/internal/logs"
	"github.com/zhlii/wechat-box/rest/internal/rpc"
)

// 创建下载文件的目录
func genDownloadDir(file string) (string, error) {
	dirs := strings.Split(file, "WeChat Files")
	if len(dirs) == 0 {
		return "", fmt.Errorf("file path does't contains WeChat Files. %s", file)
	}

	path := filepath.Join(dirs[0], "files", filepath.Base(filepath.Dir(file)))

	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return "", err
	}

	return path, nil
}

func handlerSaveMessage() {
	handlers["save_message"] = &Handler{
		Callback: func(c *rpc.Client, msg *rpc.WxMsg) {
			receiver := c.Usr.Wxid
			if msg.IsGroup {
				receiver = msg.Roomid
			}
			switch msg.Type {
			case 1:
				result := db.Db.Create(&tables.Message{
					MsgId:    msg.Id,
					Type:     1,
					IsGroup:  msg.IsGroup,
					Sender:   msg.Sender,
					Receiver: receiver,
					Content:  msg.Content,
				})
				if result.Error != nil {
					logs.Error(fmt.Sprintf("save message error. %v", result.Error))
					return
				}
			case 3:
				dir, err := genDownloadDir(msg.Extra)

				if err != nil {
					logs.Error(fmt.Sprintf("gen download dir error. filepath:%s error:%v", msg.Extra, err))
					return
				}

				if status, err := c.CmdClient.DownloadAttach(msg.Id, "", msg.Extra); status != 0 || err != nil {
					logs.Error(fmt.Sprintf("failed to download attach. status:%d error:%v", status, err))
					return
				}

				cnt := 0
				for cnt <= 10 {
					logs.Debug(fmt.Sprintf("retry %d times", cnt))

					if path, err := c.CmdClient.DecryptImage(msg.Extra, dir); err != nil {
						logs.Error(fmt.Sprintf("failed to decrypt image. path:%s error:%v", path, err))
						break
					} else {
						if path != "" {
							logs.Debug(path)

							result := db.Db.Create(&tables.Message{
								MsgId:    msg.Id,
								Type:     2,
								IsGroup:  msg.IsGroup,
								Sender:   msg.Sender,
								Receiver: receiver,
								Content:  strings.ReplaceAll(path, string(filepath.Separator), "/"),
							})
							if result.Error != nil {
								logs.Error(fmt.Sprintf("save message error. %v", result.Error))
								return
							}

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
