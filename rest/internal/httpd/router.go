package httpd

import (
	"github.com/gin-gonic/gin"
	"github.com/zhlii/wechat-box/rest/internal/rpc"
)

func Route(client *rpc.Client, rg *gin.RouterGroup) {
	ctrl := &Controller{client}

	rg.GET("login_qr", ctrl.loginQr)

	rg.GET("is_login", ctrl.isLogin)
	rg.GET("self_wxid", ctrl.getSelfWxid)
	rg.GET("self_info", ctrl.getSelfInfo)
	rg.GET("msg_types", ctrl.getMsgTypes)

	rg.GET("db_names", ctrl.getDbNames)
	rg.GET("db_tables", ctrl.getDbTables)
	rg.GET("db_query_sql", ctrl.dbSqlQuery)

	rg.GET("chatrooms", ctrl.getChatRooms)

	rg.POST("invite_chatroom_members", ctrl.inviteChatroomMembers)
	rg.POST("add_chatroom_members", ctrl.addChatRoomMembers)
	rg.DELETE("del_chatroom_members", ctrl.delChatRoomMembers)

	rg.POST("revoke_msg", ctrl.revokeMsg)
	rg.POST("forward_msg", ctrl.forwardMsg)
	rg.POST("send_txt", ctrl.sendTxt)
	rg.POST("send_img", ctrl.sendImg)
	rg.POST("send_file", ctrl.sendFile)
	rg.POST("send_rich_text", ctrl.sendRichText)
	rg.POST("send_pat_msg", ctrl.sendPatMsg)
	rg.GET("audio_msg", ctrl.getAudioMsg)
	rg.GET("ocr_result", ctrl.getOcrResult)
	rg.POST("download_image", ctrl.downloadImage)
	rg.POST("download_attach", ctrl.downloadAttach)

	rg.POST("avatars", ctrl.getAvatars)
	rg.POST("contacts", ctrl.getContacts)
	rg.POST("friends", ctrl.getFriends)
	rg.POST("user_info", ctrl.getInfoByWxid)
	rg.POST("refresh_pyq", ctrl.refreshPyq)
	rg.POST("accept_new_friend", ctrl.acceptNewFriend)
	rg.POST("receive_transfer", ctrl.receiveTransfer)
}
