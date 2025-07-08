package logic

import (
	"github.com/LagrangeDev/LagrangeGo/client"
	"github.com/LagrangeDev/LagrangeGo/client/event"
	"github.com/LagrangeDev/LagrangeGo/message"
)

// PrivateMessageHandler 私聊消息处理器
type PrivateMessageHandler struct{}

func (h *PrivateMessageHandler) Handle(client *client.QQClient, msg interface{}) error {
	if event, ok := msg.(*message.PrivateMessage); ok {
		client.SendPrivateMessage(event.Sender.Uin, []message.IMessageElement{message.NewText("Hello World!")})
	}
	return nil
}

// GroupMessageHandler 群消息处理器
type GroupMessageHandler struct{}

func (h *GroupMessageHandler) Handle(client *client.QQClient, msg interface{}) error {
	if event, ok := msg.(*message.GroupMessage); ok {
		client.SendGroupMessage(event.GroupUin, []message.IMessageElement{message.NewText("Hello World!")})
	}
	return nil
}

// FriendRequestHandler 好友请求处理器
type FriendRequestHandler struct{}

func (h *FriendRequestHandler) Handle(client *client.QQClient, msg interface{}) error {
	if event, ok := msg.(*event.NewFriendRequest); ok {
		_ = event
		// event.SourceUid
		// logrus.Println("UID" + event.SourceUid)
		// client.SetFriendRequest(true, event.SourceUid)
	}
	return nil
}

// RegisterCustomLogic 注册所有自定义逻辑
func RegisterCustomLogic() {
	// 注册私聊消息处理逻辑
	Manager.RegisterHandler("private", &PrivateMessageHandler{})

	// 注册群消息处理逻辑
	Manager.RegisterHandler("group", &GroupMessageHandler{})

	// 注册好友请求处理逻辑
	Manager.RegisterHandler("friend_request", &FriendRequestHandler{})
}
