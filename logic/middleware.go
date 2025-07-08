package logic

import (
	"fmt"
	"time"

	"github.com/LagrangeDev/LagrangeGo/client/event"
	"github.com/LagrangeDev/LagrangeGo/message"
	"github.com/sirupsen/logrus"
)

// LoggingMiddleware 日志中间件
func LoggingMiddleware() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx *MessageContext) error {
			start := time.Now()
			
			// 记录请求开始
			msgType := getMessageType(ctx.Message)
			logrus.Debugf("开始处理 %s 消息", msgType)
			
			err := next(ctx)
			
			// 记录请求结束
			duration := time.Since(start)
			if err != nil {
				logrus.Errorf("处理 %s 消息失败 (耗时: %v): %v", msgType, duration, err)
			} else {
				logrus.Debugf("处理 %s 消息成功 (耗时: %v)", msgType, duration)
			}
			
			return err
		}
	}
}

// RecoveryMiddleware 恢复中间件
func RecoveryMiddleware() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx *MessageContext) error {
			defer func() {
				if r := recover(); r != nil {
					logrus.Errorf("消息处理器发生panic: %v", r)
				}
			}()
			return next(ctx)
		}
	}
}

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware(maxRequests int, window time.Duration) Middleware {
	requests := make(map[string][]time.Time)
	
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx *MessageContext) error {
			var userID string
			
			// 获取用户ID
			if privateMsg, ok := ctx.GetPrivateMessage(); ok {
				userID = fmt.Sprintf("private_%d", privateMsg.Sender.Uin)
			} else if groupMsg, ok := ctx.GetGroupMessage(); ok {
				userID = fmt.Sprintf("group_%d_%d", groupMsg.GroupUin, groupMsg.Sender.Uin)
			} else {
				return next(ctx)
			}
			
			now := time.Now()
			
			// 清理过期请求
			userRequests := requests[userID]
			validRequests := make([]time.Time, 0)
			for _, reqTime := range userRequests {
				if now.Sub(reqTime) < window {
					validRequests = append(validRequests, reqTime)
				}
			}
			
			// 检查限流
			if len(validRequests) >= maxRequests {
				logrus.Warnf("用户 %s 触发限流", userID)
				return fmt.Errorf("请求过于频繁，请稍后再试")
			}
			
			// 记录请求
			validRequests = append(validRequests, now)
			requests[userID] = validRequests
			
			return next(ctx)
		}
	}
}

// AuthMiddleware 认证中间件
func AuthMiddleware(allowedUsers []uint32) Middleware {
	userSet := make(map[uint32]bool)
	for _, user := range allowedUsers {
		userSet[user] = true
	}
	
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx *MessageContext) error {
			var userID uint32
			
			// 获取用户ID
			if privateMsg, ok := ctx.GetPrivateMessage(); ok {
				userID = privateMsg.Sender.Uin
			} else if groupMsg, ok := ctx.GetGroupMessage(); ok {
				userID = groupMsg.Sender.Uin
			} else {
				return next(ctx)
			}
			
			// 检查权限
			if !userSet[userID] {
				logrus.Warnf("未授权用户 %d 尝试访问", userID)
				return fmt.Errorf("无权限访问")
			}
			
			ctx.Set("user_id", userID)
			ctx.Set("authorized", true)
			
			return next(ctx)
		}
	}
}

// GroupOnlyMiddleware 仅群聊中间件
func GroupOnlyMiddleware() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx *MessageContext) error {
			if _, ok := ctx.GetGroupMessage(); !ok {
				return fmt.Errorf("此命令仅在群聊中可用")
			}
			return next(ctx)
		}
	}
}

// PrivateOnlyMiddleware 仅私聊中间件
func PrivateOnlyMiddleware() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx *MessageContext) error {
			if _, ok := ctx.GetPrivateMessage(); !ok {
				return fmt.Errorf("此命令仅在私聊中可用")
			}
			return next(ctx)
		}
	}
}

// MetricsMiddleware 指标中间件
func MetricsMiddleware() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx *MessageContext) error {
			msgType := getMessageType(ctx.Message)
			
			// 记录指标
			ctx.Set("message_type", msgType)
			ctx.Set("processed_at", time.Now())
			
			err := next(ctx)
			
			// 记录处理结果
			if err != nil {
				ctx.Set("processing_error", err.Error())
			}
			
			return err
		}
	}
}

// getMessageType 获取消息类型
func getMessageType(msg interface{}) string {
	switch msg.(type) {
	case *message.PrivateMessage:
		return "private"
	case *message.GroupMessage:
		return "group"
	case *event.NewFriendRequest:
		return "friend_request"
	default:
		return "unknown"
	}
}

// ConditionalMiddleware 条件中间件
func ConditionalMiddleware(condition func(*MessageContext) bool, middleware Middleware) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx *MessageContext) error {
			if condition(ctx) {
				return middleware(next)(ctx)
			}
			return next(ctx)
		}
	}
}

// ChainMiddleware 链式中间件
func ChainMiddleware(middlewares ...Middleware) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}