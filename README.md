# LagrangeGo-Template

A template for LagrangeGo

[![Go Report Card](https://goreportcard.com/badge/github.com/ExquisiteCore/LagrangeGo-Template)](https://goreportcard.com/report/github.com/ExquisiteCore/LagrangeGo-Template)

基于 [LagrangeGo](https://github.com/LagrangeDev/LagrangeGo) 的Bot 模板参考自[MiraiGo-Template](https://github.com/Logiase/MiraiGo-Template)

## 项目特点

- 使用**组合优于继承**的设计原则
- 实现**依赖注入**模式，提高代码可测试性
- 清晰的**分层架构**，职责分离
- 统一的**消息处理器接口**，易于扩展

## 项目结构

```
├── app/           # 应用层，依赖注入容器
├── bot/           # Bot层，包含登录和认证管理
├── config/        # 配置层
├── logic/         # 逻辑层，消息处理器
├── utils/         # 工具层
├── main.go        # 主程序入口
└── application.toml # 配置文件
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

克隆本项目

在[logic/custom_logic.go](./logic/custom_logic.go)注册逻辑

```go
// RegisterCustomLogic 注册所有自定义逻辑
func RegisterCustomLogic() {
    // 注册私聊消息处理逻辑
    Manager.RegisterHandler("private", &PrivateMessageHandler{})
    
    // 注册群消息处理逻辑
    Manager.RegisterHandler("group", &GroupMessageHandler{})
    
    // 注册好友请求处理逻辑
    Manager.RegisterHandler("friend_request", &FriendRequestHandler{})
}
```

## 消息处理器实现

所有消息处理器都需要实现 `MessageHandler` 接口：

```go
type MessageHandler interface {
    Handle(client *client.QQClient, msg interface{}) error
}
```

示例实现：

```go
// PrivateMessageHandler 私聊消息处理器
type PrivateMessageHandler struct{}

func (h *PrivateMessageHandler) Handle(client *client.QQClient, msg interface{}) error {
    if event, ok := msg.(*message.PrivateMessage); ok {
        client.SendPrivateMessage(event.Sender.Uin, []message.IMessageElement{message.NewText("Hello World!")})
    }
    return nil
}
```

## 架构设计

### 组合优于继承

- `Bot` 结构体通过组合 `*client.QQClient` 而非继承
- `LoginManager` 和 `AuthManager` 各自负责特定功能
- `LogicManager` 通过组合管理各种处理器

### 依赖注入

使用 `app.Container` 统一管理所有依赖：

```go
container := app.NewContainer()
err := container.Initialize()
bot := container.GetBot()
logicManager := container.GetLogicManager()
```

### 接口设计

统一的 `MessageHandler` 接口使得处理器易于扩展和测试，支持多种消息类型的处理。

## 引入的第三方 go module

- [LagrangeGo](https://github.com/LagrangeDev/LagrangeGo)
    核心协议库
- [toml](https://github.com/BurntSushi/toml)
    用于解析配置文件
- [logrus](https://github.com/sirupsen/logrus)
    功能丰富的Logger
- [qrterminal](https://github.com/mdp/qrterminal)
    终端二维码显示
- [qrcode](https://github.com/tuotoo/qrcode)
    二维码解析
