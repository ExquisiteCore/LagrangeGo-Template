package logic

import (
	"github.com/ExquisiteCore/LagrangeGo-Template/bot"
	"github.com/LagrangeDev/LagrangeGo/client"
	"github.com/LagrangeDev/LagrangeGo/client/event"
	"github.com/LagrangeDev/LagrangeGo/message"
)

// 定义不同类型事件的处理函数类型
type PrivateMessageHandler func(*client.QQClient, *message.PrivateMessage)
type GroupMessageHandler func(*client.QQClient, *message.GroupMessage)
type NewFriendRequestHandler func(*client.QQClient, *event.NewFriendRequest)

// LogicManager 管理所有自定义逻辑
type LogicManager struct {
	privateMessageHandlers   []PrivateMessageHandler
	groupMessageHandlers     []GroupMessageHandler
	newFriendRequestHandlers []NewFriendRequestHandler
}

// 全局 LogicManager 实例
var Manager = &LogicManager{}

// RegisterPrivateMessageHandler 注册私聊消息处理函数
func (lm *LogicManager) RegisterPrivateMessageHandler(handler PrivateMessageHandler) {
	lm.privateMessageHandlers = append(lm.privateMessageHandlers, handler)
}

// RegisterGroupMessageHandler 注册群消息处理函数
func (lm *LogicManager) RegisterGroupMessageHandler(handler GroupMessageHandler) {
	lm.groupMessageHandlers = append(lm.groupMessageHandlers, handler)
}

// RegisterNewFriendRequestHandler 注册新朋友请求处理函数
func (lm *LogicManager) RegisterNewFriendRequestHandler(handler NewFriendRequestHandler) {
	lm.newFriendRequestHandlers = append(lm.newFriendRequestHandlers, handler)
}

// SetupLogic 设置所有逻辑处理
func SetupLogic() {
	bot.QQClient.PrivateMessageEvent.Subscribe(func(client *client.QQClient, event *message.PrivateMessage) {
		for _, handler := range Manager.privateMessageHandlers {
			handler(client, event)
		}
	})

	bot.QQClient.GroupMessageEvent.Subscribe(func(client *client.QQClient, event *message.GroupMessage) {
		for _, handler := range Manager.groupMessageHandlers {
			handler(client, event)
		}
	})

	bot.QQClient.NewFriendRequestEvent.Subscribe(func(client *client.QQClient, event *event.NewFriendRequest) {
		for _, handler := range Manager.newFriendRequestHandlers {
			handler(client, event)
		}
	})
}
