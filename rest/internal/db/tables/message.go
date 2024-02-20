package tables

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	MsgId    uint64 `gorm:"uniqueIndex"` // 消息 id
	Type     uint32 // 消息类型
	IsGroup  bool   // 是否群消息
	Sender   string `gorm:"size:50"`  // 消息发送者
	Receiver string `gorm:"size:50"`  // 消息接收者
	Content  string `gorm:"size:500"` // 消息内容
}
