package logic

import (
	"github.com/LagrangeDev/LagrangeGo/client"
	"github.com/LagrangeDev/LagrangeGo/client/event"
	"github.com/LagrangeDev/LagrangeGo/message"
)

// MessageHandler 通用消息处理器接口
type MessageHandler interface {
	Handle(client *client.QQClient, msg interface{}) error
}

// EventHandlers 事件处理器容器
type EventHandlers struct {
	privateHandlers       []MessageHandler
	groupHandlers         []MessageHandler
	friendRequestHandlers []MessageHandler
}

// NewEventHandlers 创建新的事件处理器容器
func NewEventHandlers() *EventHandlers {
	return &EventHandlers{
		privateHandlers:       make([]MessageHandler, 0),
		groupHandlers:         make([]MessageHandler, 0),
		friendRequestHandlers: make([]MessageHandler, 0),
	}
}

// RegisterPrivateHandler 注册私聊消息处理器
func (eh *EventHandlers) RegisterPrivateHandler(handler MessageHandler) {
	eh.privateHandlers = append(eh.privateHandlers, handler)
}

// RegisterGroupHandler 注册群消息处理器
func (eh *EventHandlers) RegisterGroupHandler(handler MessageHandler) {
	eh.groupHandlers = append(eh.groupHandlers, handler)
}

// RegisterFriendRequestHandler 注册好友请求处理器
func (eh *EventHandlers) RegisterFriendRequestHandler(handler MessageHandler) {
	eh.friendRequestHandlers = append(eh.friendRequestHandlers, handler)
}

// LogicManager 逻辑管理器
type LogicManager struct {
	client   *client.QQClient
	handlers *EventHandlers
}

// NewLogicManager 创建新的逻辑管理器
func NewLogicManager(client *client.QQClient) *LogicManager {
	return &LogicManager{
		client:   client,
		handlers: NewEventHandlers(),
	}
}

// RegisterHandler 注册处理器
func (lm *LogicManager) RegisterHandler(handlerType string, handler MessageHandler) {
	switch handlerType {
	case "private":
		lm.handlers.RegisterPrivateHandler(handler)
	case "group":
		lm.handlers.RegisterGroupHandler(handler)
	case "friend_request":
		lm.handlers.RegisterFriendRequestHandler(handler)
	}
}

// SetupEventListeners 设置事件监听器
func (lm *LogicManager) SetupEventListeners() {
	lm.client.PrivateMessageEvent.Subscribe(func(client *client.QQClient, event *message.PrivateMessage) {
		for _, handler := range lm.handlers.privateHandlers {
			handler.Handle(client, event)
		}
	})

	lm.client.GroupMessageEvent.Subscribe(func(client *client.QQClient, event *message.GroupMessage) {
		for _, handler := range lm.handlers.groupHandlers {
			handler.Handle(client, event)
		}
	})

	lm.client.NewFriendRequestEvent.Subscribe(func(client *client.QQClient, event *event.NewFriendRequest) {
		for _, handler := range lm.handlers.friendRequestHandlers {
			handler.Handle(client, event)
		}
	})
}

// 全局 LogicManager 实例
var Manager *LogicManager

// SetupLogic 设置逻辑处理
func SetupLogic(client *client.QQClient) {
	Manager = NewLogicManager(client)
	Manager.SetupEventListeners()
}
