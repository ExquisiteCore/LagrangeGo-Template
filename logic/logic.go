package logic

import (
	"github.com/ExquisiteCore/LagrangeGo-Template/utils"
	"github.com/LagrangeDev/LagrangeGo/client"
	"github.com/LagrangeDev/LagrangeGo/client/event"
	"github.com/LagrangeDev/LagrangeGo/message"
)

// LogicManager 新的逻辑管理器
type LogicManager struct {
	client   *client.QQClient
	router   *Router
	eventBus *EventBus
}

// NewLogicManager 创建新的逻辑管理器
func NewLogicManager(client *client.QQClient) *LogicManager {
	return &LogicManager{
		client:   client,
		router:   NewRouter(),
		eventBus: NewEventBus(),
	}
}

// GetRouter 获取路由器
func (lm *LogicManager) GetRouter() *Router {
	return lm.router
}

// GetEventBus 获取事件总线
func (lm *LogicManager) GetEventBus() *EventBus {
	return lm.eventBus
}

// UseMiddleware 使用全局中间件
func (lm *LogicManager) UseMiddleware(middleware Middleware) {
	lm.router.Use(middleware)
}

// AddRoute 添加路由
func (lm *LogicManager) AddRoute(route *Route) {
	lm.router.AddRoute(route)
}

// HandlePrivateMessage 处理私聊消息的便捷方法
func (lm *LogicManager) HandlePrivateMessage(handler HandlerFunc, matchers ...Matcher) {
	route := NewRoute("private_message", NewHandlerAdapter(handler))
	route.Match(NewMessageTypeMatcher("private"))
	for _, matcher := range matchers {
		route.Match(matcher)
	}
	lm.AddRoute(route)
}

// HandleGroupMessage 处理群消息的便捷方法
func (lm *LogicManager) HandleGroupMessage(handler HandlerFunc, matchers ...Matcher) {
	route := NewRoute("group_message", NewHandlerAdapter(handler))
	route.Match(NewMessageTypeMatcher("group"))
	for _, matcher := range matchers {
		route.Match(matcher)
	}
	lm.AddRoute(route)
}

// HandleFriendRequest 处理好友请求的便捷方法
func (lm *LogicManager) HandleFriendRequest(handler HandlerFunc, matchers ...Matcher) {
	route := NewRoute("friend_request", NewHandlerAdapter(handler))
	route.Match(NewMessageTypeMatcher("friend_request"))
	for _, matcher := range matchers {
		route.Match(matcher)
	}
	lm.AddRoute(route)
}

// HandleCommand 处理命令的便捷方法
func (lm *LogicManager) HandleCommand(prefix string, command string, handler HandlerFunc, middlewares ...Middleware) {
	route := NewRoute("command_"+command, NewHandlerAdapter(handler))
	route.Match(NewCommandMatcher(prefix, command))
	for _, middleware := range middlewares {
		route.Use(middleware)
	}
	lm.AddRoute(route)
}

// SetupEventListeners 设置事件监听器
func (lm *LogicManager) SetupEventListeners() {
	// 私聊消息事件
	lm.client.PrivateMessageEvent.Subscribe(func(client *client.QQClient, event *message.PrivateMessage) {
		ctx := NewMessageContext(client, event)
		lm.processMessage(ctx)
	})

	// 群消息事件
	lm.client.GroupMessageEvent.Subscribe(func(client *client.QQClient, event *message.GroupMessage) {
		ctx := NewMessageContext(client, event)
		lm.processMessage(ctx)
	})

	// 好友请求事件
	lm.client.NewFriendRequestEvent.Subscribe(func(client *client.QQClient, event *event.NewFriendRequest) {
		ctx := NewMessageContext(client, event)
		lm.processMessage(ctx)
	})
}

// processMessage 处理消息
func (lm *LogicManager) processMessage(ctx *MessageContext) {
	// 发布消息接收事件
	PublishMessageReceived(ctx)
	
	// 通过路由器处理消息
	lm.router.Handle(ctx)
	
	// 发布消息处理完成事件
	PublishMessageProcessed(ctx)
}

// Close 关闭逻辑管理器
func (lm *LogicManager) Close() {
	lm.eventBus.Close()
}

// 全局 LogicManager 实例
var Manager *LogicManager

// SetupLogic 设置逻辑处理
func SetupLogic(client *client.QQClient) {
	Manager = NewLogicManager(client)
	
	// 设置默认中间件
	Manager.UseMiddleware(RecoveryMiddleware())
	Manager.UseMiddleware(LoggingMiddleware())
	Manager.UseMiddleware(MetricsMiddleware())
	
	// 设置事件监听器
	Manager.SetupEventListeners()
	
	utils.Info("逻辑管理器初始化完成")
}

// 向后兼容的方法
type MessageHandler interface {
	Handle(client *client.QQClient, msg interface{}) error
}

type HandlerWrapper struct {
	handler MessageHandler
}

func (hw *HandlerWrapper) Handle(ctx *MessageContext) error {
	return hw.handler.Handle(ctx.Client, ctx.Message)
}

// RegisterHandler 注册处理器（向后兼容）
func (lm *LogicManager) RegisterHandler(handlerType string, handler MessageHandler) {
	wrapper := &HandlerWrapper{handler: handler}
	route := NewRoute(handlerType+"_handler", wrapper)
	route.Match(NewMessageTypeMatcher(handlerType))
	lm.AddRoute(route)
}
