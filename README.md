# LagrangeGo-Template

A template for LagrangeGo

[![Go Report Card](https://goreportcard.com/badge/github.com/ExquisiteCore/LagrangeGo-Template)](https://goreportcard.com/report/github.com/ExquisiteCore/LagrangeGo-Template)

基于 [LagrangeGo](https://github.com/LagrangeDev/LagrangeGo) 的Bot 模板参考自[MiraiGo-Template](https://github.com/Logiase/MiraiGo-Template)

## 项目特点

- 使用**组合优于继承**的设计原则
- 实现**依赖注入**模式，提高代码可测试性
- 清晰的**分层架构**，职责分离
- 统一的**消息处理器接口**，易于扩展
- **中间件支持**，横切关注点统一处理
- **事件总线**，松耦合的组件通信
- **强大的路由系统**，支持多种匹配条件
- **统一的日志系统**，支持文件输出和彩色终端

## 项目结构

```
├── app/              # 应用层，依赖注入容器
├── bot/              # Bot层，登录、认证、连接管理
│   ├── auth.go       # 认证管理器
│   ├── bot.go        # Bot主结构
│   ├── connection.go # 连接管理器
│   ├── login.go      # 登录策略
│   └── qrcode.go     # 二维码处理器
├── config/           # 配置层
├── logic/            # 逻辑层，消息处理器
│   ├── eventbus.go   # 事件总线
│   ├── handlers.go   # 示例处理器
│   ├── logic.go      # 逻辑管理器
│   ├── matcher.go    # 消息匹配器
│   ├── middleware.go # 中间件
│   └── router.go     # 路由系统
├── utils/            # 工具层
│   └── log.go        # 统一日志系统
├── main.go           # 主程序入口
└── application.toml  # 配置文件
```

## 基础配置

账号配置[application.toml](./application.toml)

```toml
[bot]
# 账号 必填
account = 114514
# 密码 选填
password = "pwd"
# 签名服务器 选填
signServer = "https://sign.lagrangecore.org/api/sign/25765"
```

不配置密码的话将使用扫码登录

## 快速入门

### 1. 克隆项目

```bash
git clone https://github.com/ExquisiteCore/LagrangeGo-Template.git
cd LagrangeGo-Template
```

### 2. 配置账号

编辑 `application.toml` 文件，配置你的QQ账号。

### 3. 注册消息处理器

在 `logic/handlers.go` 中注册你的消息处理逻辑：

```go
// RegisterCustomLogic 注册所有自定义逻辑
func RegisterCustomLogic() {
    // 注册ping命令
    Manager.HandleCommand("/", "ping", func(ctx *MessageContext) error {
        // 处理ping命令
        return nil
    })
    
    // 注册文本匹配处理器
    Manager.HandleGroupMessage(func(ctx *MessageContext) error {
        text := ctx.GetMessageText()
        if text == "hello" {
            // 处理hello消息
        }
        return nil
    }, NewTextMatcher("hello", false))
}
```

### 4. 运行项目

```bash
go run main.go
```

## 消息处理器

### 基础处理器

所有消息处理器都需要实现 `MessageHandler` 接口：

```go
type MessageHandler interface {
    Handle(ctx *MessageContext) error
}
```

### 便捷方法

LogicManager 提供了多种便捷方法：

```go
// 处理命令
Manager.HandleCommand("/", "ping", handlerFunc)

// 处理私聊消息
Manager.HandlePrivateMessage(handlerFunc, matchers...)

// 处理群消息
Manager.HandleGroupMessage(handlerFunc, matchers...)

// 处理好友请求
Manager.HandleFriendRequest(handlerFunc, matchers...)
```

### 消息匹配器

支持多种匹配条件：

```go
// 文本匹配
NewTextMatcher("hello", false)

// 正则匹配
NewRegexMatcher(`^\d+$`)

// 前缀匹配
NewPrefixMatcher("/")

// 发送者匹配
NewSenderMatcher(12345, 67890)

// 群聊匹配
NewGroupMatcher(111111, 222222)

// 命令匹配
NewCommandMatcher("/", "ping", "help")

// 组合匹配
NewAndMatcher(matcher1, matcher2)
NewOrMatcher(matcher1, matcher2)
```

## 中间件系统

### 内置中间件

- **LoggingMiddleware**: 请求日志记录
- **RecoveryMiddleware**: 错误恢复
- **RateLimitMiddleware**: 频率限制
- **AuthMiddleware**: 权限验证
- **MetricsMiddleware**: 指标收集

### 使用示例

```go
// 全局中间件
Manager.UseMiddleware(LoggingMiddleware())
Manager.UseMiddleware(RecoveryMiddleware())

// 路由中间件
Manager.HandleCommand("/", "admin", handlerFunc, 
    AuthMiddleware([]uint32{12345}),
    RateLimitMiddleware(5, time.Minute),
)
```

## 事件系统

### 事件总线

项目内置事件总线，支持异步事件处理：

```go
// 订阅事件
Manager.GetEventBus().Subscribe(EventTypeCommandExecuted, func(ctx context.Context, event Event) error {
    // 处理命令执行事件
    return nil
})

// 发布事件
PublishMessageReceived(msgContext)
PublishCommandExecuted(msgContext, "ping")
```

### 预定义事件

- `EventTypeMessageReceived`: 消息接收事件
- `EventTypeMessageProcessed`: 消息处理完成事件
- `EventTypeCommandExecuted`: 命令执行事件
- `EventTypeError`: 错误事件

## 日志系统

### 日志配置

```go
config := &utils.LogConfig{
    Level:       utils.InfoLevel,
    EnableFile:  true,
    EnableColor: true,
    LogDir:      "logs",
    LogFile:     "bot.log",
    MaxSize:     10, // MB
    MaxBackups:  5,
    MaxAge:      30, // days
    Format:      "text", // text or json
}
utils.InitWithConfig(config)
```

### 日志使用

```go
// 基础日志
utils.Info("这是一条信息")
utils.Infof("用户 %d 发送了消息", userID)

// 带字段的日志
utils.WithField("user_id", 12345).Info("用户操作")
utils.WithFields(map[string]interface{}{
    "user_id": 12345,
    "action":  "login",
}).Info("用户登录")
```

## Bot架构

### 登录策略

支持多种登录方式：

- **FastLoginStrategy**: 快速登录（使用保存的签名）
- **QRCodeLoginStrategy**: 二维码登录

### 连接管理

- 自动重连机制
- 连接状态监控
- 心跳包发送
- 连接事件处理

### 认证管理

- 签名文件管理
- 自动加载和保存
- 签名有效性检查

## 依赖注入

使用依赖注入容器管理所有组件：

```go
container := app.NewContainer()
err := container.Initialize()

bot := container.GetBot()
logicManager := container.GetLogicManager()
client := container.GetClient()
config := container.GetConfig()
```

## 开发指南

### 添加新的消息处理器

1. 创建处理器结构体
2. 实现 `MessageHandler` 接口
3. 在 `RegisterCustomLogic` 中注册

### 添加新的中间件

1. 实现 `Middleware` 函数类型
2. 在需要的路由上使用

### 添加新的匹配器

1. 实现 `Matcher` 接口
2. 在路由中使用

### 添加新的事件

1. 定义事件类型常量
2. 创建事件结构体
3. 发布和订阅事件

## 测试和部署

### 编译

```bash
go build -o bot .
```

### 运行

```bash
./bot
```

### Docker部署

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o bot .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/bot .
COPY --from=builder /app/application.toml .
CMD ["./bot"]
```

## 引入的第三方库

- [LagrangeGo](https://github.com/LagrangeDev/LagrangeGo) - 核心协议库
- [logrus](https://github.com/sirupsen/logrus) - 日志库
- [toml](https://github.com/BurntSushi/toml) - 配置文件解析
- [qrterminal](https://github.com/mdp/qrterminal) - 终端二维码显示
- [qrcode](https://github.com/tuotoo/qrcode) - 二维码解析
- [colorable](https://github.com/mattn/go-colorable) - 跨平台彩色终端

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！

## 更新日志

### v2.0.0
- 重构整个架构，采用组合优于继承
- 实现依赖注入模式
- 新增强大的路由系统和中间件支持
- 添加事件总线机制
- 统一日志系统
- 完善的错误处理和重连机制

### v1.0.0
- 初始版本
- 基础消息处理功能
