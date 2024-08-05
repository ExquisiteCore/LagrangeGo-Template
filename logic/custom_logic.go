package logic

import (
	"github.com/LagrangeDev/LagrangeGo/client"
	"github.com/LagrangeDev/LagrangeGo/message"
)

// RegisterCustomLogic 注册所有自定义逻辑
func RegisterCustomLogic() {
	// 注册私聊消息处理逻辑
	Manager.RegisterPrivateMessageHandler(func(client *client.QQClient, event *message.PrivateMessage) {
		client.SendPrivateMessage(event.Sender.Uin, []message.IMessageElement{message.NewText("Hello World!")})
	})

	// 注册群消息处理逻辑
	Manager.RegisterGroupMessageHandler(func(client *client.QQClient, event *message.GroupMessage) {
		client.SendGroupMessage(event.GroupUin, []message.IMessageElement{message.NewText("Hello World!")})
	})
}

//
