package logic

import (
	"context"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// Event 事件接口
type Event interface {
	GetType() string
	GetData() interface{}
	GetTimestamp() time.Time
}

// BaseEvent 基础事件结构
type BaseEvent struct {
	Type      string
	Data      interface{}
	Timestamp time.Time
}

func (e *BaseEvent) GetType() string {
	return e.Type
}

func (e *BaseEvent) GetData() interface{} {
	return e.Data
}

func (e *BaseEvent) GetTimestamp() time.Time {
	return e.Timestamp
}

// NewEvent 创建新事件
func NewEvent(eventType string, data interface{}) *BaseEvent {
	return &BaseEvent{
		Type:      eventType,
		Data:      data,
		Timestamp: time.Now(),
	}
}

// EventHandler 事件处理器
type EventHandler func(ctx context.Context, event Event) error

// EventBus 事件总线
type EventBus struct {
	handlers map[string][]EventHandler
	mu       sync.RWMutex
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

// NewEventBus 创建新的事件总线
func NewEventBus() *EventBus {
	ctx, cancel := context.WithCancel(context.Background())
	return &EventBus{
		handlers: make(map[string][]EventHandler),
		ctx:      ctx,
		cancel:   cancel,
	}
}

// Subscribe 订阅事件
func (bus *EventBus) Subscribe(eventType string, handler EventHandler) {
	bus.mu.Lock()
	defer bus.mu.Unlock()
	
	if bus.handlers[eventType] == nil {
		bus.handlers[eventType] = make([]EventHandler, 0)
	}
	
	bus.handlers[eventType] = append(bus.handlers[eventType], handler)
	logrus.Debugf("订阅事件类型: %s", eventType)
}

// Unsubscribe 取消订阅事件
func (bus *EventBus) Unsubscribe(eventType string, handler EventHandler) {
	bus.mu.Lock()
	defer bus.mu.Unlock()
	
	handlers := bus.handlers[eventType]
	if handlers == nil {
		return
	}
	
	// 移除指定的处理器
	for i, h := range handlers {
		if &h == &handler {
			bus.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}
	
	logrus.Debugf("取消订阅事件类型: %s", eventType)
}

// Publish 发布事件
func (bus *EventBus) Publish(event Event) {
	bus.mu.RLock()
	handlers := make([]EventHandler, len(bus.handlers[event.GetType()]))
	copy(handlers, bus.handlers[event.GetType()])
	bus.mu.RUnlock()
	
	if len(handlers) == 0 {
		return
	}
	
	logrus.Debugf("发布事件: %s", event.GetType())
	
	// 异步处理事件
	for _, handler := range handlers {
		bus.wg.Add(1)
		go func(h EventHandler) {
			defer bus.wg.Done()
			defer func() {
				if r := recover(); r != nil {
					logrus.Errorf("事件处理器发生panic: %v", r)
				}
			}()
			
			ctx, cancel := context.WithTimeout(bus.ctx, 30*time.Second)
			defer cancel()
			
			if err := h(ctx, event); err != nil {
				logrus.Errorf("处理事件 %s 时发生错误: %v", event.GetType(), err)
			}
		}(handler)
	}
}

// PublishSync 同步发布事件
func (bus *EventBus) PublishSync(event Event) error {
	bus.mu.RLock()
	handlers := make([]EventHandler, len(bus.handlers[event.GetType()]))
	copy(handlers, bus.handlers[event.GetType()])
	bus.mu.RUnlock()
	
	if len(handlers) == 0 {
		return nil
	}
	
	logrus.Debugf("同步发布事件: %s", event.GetType())
	
	ctx, cancel := context.WithTimeout(bus.ctx, 30*time.Second)
	defer cancel()
	
	for _, handler := range handlers {
		if err := handler(ctx, event); err != nil {
			logrus.Errorf("处理事件 %s 时发生错误: %v", event.GetType(), err)
			return err
		}
	}
	
	return nil
}

// Close 关闭事件总线
func (bus *EventBus) Close() {
	bus.cancel()
	bus.wg.Wait()
	logrus.Info("事件总线已关闭")
}

// GetSubscriberCount 获取订阅者数量
func (bus *EventBus) GetSubscriberCount(eventType string) int {
	bus.mu.RLock()
	defer bus.mu.RUnlock()
	return len(bus.handlers[eventType])
}

// GetAllEventTypes 获取所有事件类型
func (bus *EventBus) GetAllEventTypes() []string {
	bus.mu.RLock()
	defer bus.mu.RUnlock()
	
	types := make([]string, 0, len(bus.handlers))
	for eventType := range bus.handlers {
		types = append(types, eventType)
	}
	return types
}

// MessageEvent 消息事件
type MessageEvent struct {
	*BaseEvent
	MessageContext *MessageContext
}

// NewMessageEvent 创建消息事件
func NewMessageEvent(eventType string, ctx *MessageContext) *MessageEvent {
	return &MessageEvent{
		BaseEvent:      NewEvent(eventType, ctx.Message),
		MessageContext: ctx,
	}
}

// 预定义事件类型
const (
	EventTypeMessageReceived  = "message.received"
	EventTypeMessageProcessed = "message.processed"
	EventTypeCommandExecuted  = "command.executed"
	EventTypeUserJoined       = "user.joined"
	EventTypeUserLeft         = "user.left"
	EventTypeError            = "error.occurred"
)

// 全局事件总线实例
var GlobalEventBus = NewEventBus()

// PublishMessageReceived 发布消息接收事件
func PublishMessageReceived(ctx *MessageContext) {
	event := NewMessageEvent(EventTypeMessageReceived, ctx)
	GlobalEventBus.Publish(event)
}

// PublishMessageProcessed 发布消息处理完成事件
func PublishMessageProcessed(ctx *MessageContext) {
	event := NewMessageEvent(EventTypeMessageProcessed, ctx)
	GlobalEventBus.Publish(event)
}

// PublishCommandExecuted 发布命令执行事件
func PublishCommandExecuted(ctx *MessageContext, command string) {
	ctx.Set("executed_command", command)
	event := NewMessageEvent(EventTypeCommandExecuted, ctx)
	GlobalEventBus.Publish(event)
}

// PublishError 发布错误事件
func PublishError(err error, ctx *MessageContext) {
	data := map[string]interface{}{
		"error":   err.Error(),
		"context": ctx,
	}
	event := NewEvent(EventTypeError, data)
	GlobalEventBus.Publish(event)
}