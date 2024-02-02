package callback

import (
	"encoding/xml"
	"fmt"
	"regexp"

	"github.com/zhlii/wechat-box/rest/internal/logs"
	"github.com/zhlii/wechat-box/rest/internal/rpc"
)

func handlerAutoAcceptFriendInvite() {
	handlers["auto_accept_friend_invite"] = &Handler{
		Callback: func(c *rpc.Client, msg *rpc.WxMsg) {
			switch msg.Type {
			case 37:
				type MSG struct {
					Encryptusername string `xml:"encryptusername,attr"`
					Ticket          string `xml:"ticket,attr"`
					Scene           int32  `xml:"scene,attr"`
				}

				var result MSG
				err := xml.Unmarshal([]byte(msg.Content), &result)
				if err != nil {
					logs.Error(fmt.Sprintf("Error decoding xml: %s. error:%v\n", msg.Content, err))
					return
				}

				_, err = c.CmdClient.AcceptNewFriend(result.Encryptusername, result.Ticket, result.Scene)
				if err != nil {
					logs.Error(fmt.Sprintf("AcceptFriend error: %v\n", err))
				}
			case 10000:
				re := regexp.MustCompile("你已添加了(.*?)，现在可以开始聊天了。")

				matches := re.FindStringSubmatch(msg.Content)
				if len(matches) > 1 {

					c.FreshContacts()

					nickName := matches[1]
					logs.Debug(fmt.Sprintf("you and %s is friend now.", nickName))

					c.CmdClient.SendTxt(fmt.Sprintf("Hi %s", nickName), msg.Sender, "")
				}
			}
		},
	}
}
