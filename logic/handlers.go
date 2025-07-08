package logic

import (
	"context"
	
	"github.com/ExquisiteCore/LagrangeGo-Template/utils"
	"github.com/LagrangeDev/LagrangeGo/client"
	"github.com/LagrangeDev/LagrangeGo/client/event"
	"github.com/LagrangeDev/LagrangeGo/message"
)

// EchoHandler 回声处理器
type EchoHandler struct{}

func (h *EchoHandler) Handle(ctx *MessageContext) error {
	text := ctx.GetMessageText()
	if text == "" {
		return nil
	}
	
	if privateMsg, ok := ctx.GetPrivateMessage(); ok {
		ctx.Client.SendPrivateMessage(privateMsg.Sender.Uin, []message.IMessageElement{
			message.NewText("你说了: " + text),
		})
	} else if groupMsg, ok := ctx.GetGroupMessage(); ok {
		ctx.Client.SendGroupMessage(groupMsg.GroupUin, []message.IMessageElement{
			message.NewText("你说了: " + text),
		})
	}
	
	return nil
}

// PingHandler ping命令处理器
type PingHandler struct{}

func (h *PingHandler) Handle(ctx *MessageContext) error {
	response := "pong!"
	
	if privateMsg, ok := ctx.GetPrivateMessage(); ok {
		ctx.Client.SendPrivateMessage(privateMsg.Sender.Uin, []message.IMessageElement{
			message.NewText(response),
		})
	} else if groupMsg, ok := ctx.GetGroupMessage(); ok {
		ctx.Client.SendGroupMessage(groupMsg.GroupUin, []message.IMessageElement{
			message.NewText(response),
		})
	}
	
	return nil
}

// HelpHandler 帮助命令处理器
type HelpHandler struct{}

func (h *HelpHandler) Handle(ctx *MessageContext) error {
	help := `可用命令:
/ping - 测试连接
/help - 显示帮助
/echo <消息> - 回声消息`
	
	if privateMsg, ok := ctx.GetPrivateMessage(); ok {
		ctx.Client.SendPrivateMessage(privateMsg.Sender.Uin, []message.IMessageElement{
			message.NewText(help),
		})
	} else if groupMsg, ok := ctx.GetGroupMessage(); ok {
		ctx.Client.SendGroupMessage(groupMsg.GroupUin, []message.IMessageElement{
			message.NewText(help),
		})
	}
	
	return nil
}

// FriendRequestHandler 好友请求处理器
type FriendRequestHandler struct{}

func (h *FriendRequestHandler) Handle(ctx *MessageContext) error {
	if friendReq, ok := ctx.GetFriendRequest(); ok {
		utils.Infof("收到好友请求: %s", friendReq.SourceUID)
		// 可以在这里实现自动接受好友请求的逻辑
		// ctx.Client.SetFriendRequest(true, friendReq.SourceUID)
	}
	return nil
}

// RegisterCustomLogic 注册所有自定义逻辑
func RegisterCustomLogic() {
	if Manager == nil {
		utils.Error("LogicManager 未初始化")
		return
	}
	
	// 注册ping命令
	Manager.HandleCommand("/", "ping", func(ctx *MessageContext) error {
		handler := &PingHandler{}
		return handler.Handle(ctx)
	})
	
	// 注册help命令
	Manager.HandleCommand("/", "help", func(ctx *MessageContext) error {
		handler := &HelpHandler{}
		return handler.Handle(ctx)
	})
	
	// 注册echo命令
	Manager.HandleCommand("/", "echo", func(ctx *MessageContext) error {
		handler := &EchoHandler{}
		return handler.Handle(ctx)
	})
	
	// 注册好友请求处理
	Manager.HandleFriendRequest(func(ctx *MessageContext) error {
		handler := &FriendRequestHandler{}
		return handler.Handle(ctx)
	})
	
	// 注册基于文本匹配的处理器
	Manager.HandleGroupMessage(func(ctx *MessageContext) error {
		text := ctx.GetMessageText()
		if text == "hello" {
			if groupMsg, ok := ctx.GetGroupMessage(); ok {
				ctx.Client.SendGroupMessage(groupMsg.GroupUin, []message.IMessageElement{
					message.NewText("Hello World!"),
				})
			}
		}
		return nil
	}, NewTextMatcher("hello", false))
	
	// 注册私聊消息处理
	Manager.HandlePrivateMessage(func(ctx *MessageContext) error {
		text := ctx.GetMessageText()
		if text == "hello" {
			if privateMsg, ok := ctx.GetPrivateMessage(); ok {
				ctx.Client.SendPrivateMessage(privateMsg.Sender.Uin, []message.IMessageElement{
					message.NewText("Hello! 这是私聊回复"),
				})
			}
		}
		return nil
	}, NewTextMatcher("hello", false))
	
	// 注册事件监听器
	Manager.GetEventBus().Subscribe(EventTypeCommandExecuted, func(ctx context.Context, event Event) error {
		if msgEvent, ok := event.(*MessageEvent); ok {
			command := msgEvent.MessageContext.GetString("executed_command")
			utils.Infof("命令 %s 已执行", command)
		}
		return nil
	})
	
	utils.Info("自定义逻辑注册完成")
}

// 向后兼容的处理器实现
type PrivateMessageHandler struct{}

func (h *PrivateMessageHandler) Handle(client *client.QQClient, msg interface{}) error {
	if event, ok := msg.(*message.PrivateMessage); ok {
		client.SendPrivateMessage(event.Sender.Uin, []message.IMessageElement{
			message.NewText("Hello World!"),
		})
	}
	return nil
}

type GroupMessageHandler struct{}

func (h *GroupMessageHandler) Handle(client *client.QQClient, msg interface{}) error {
	if event, ok := msg.(*message.GroupMessage); ok {
		client.SendGroupMessage(event.GroupUin, []message.IMessageElement{
			message.NewText("Hello World!"),
		})
	}
	return nil
}

type FriendRequestHandlerOld struct{}

func (h *FriendRequestHandlerOld) Handle(client *client.QQClient, msg interface{}) error {
	if event, ok := msg.(*event.NewFriendRequest); ok {
		_ = event
		// event.SourceUid
		// logrus.Println("UID" + event.SourceUid)
		// ctx.Client.SetFriendRequest(true, event.SourceUID)
	}
	return nil
}