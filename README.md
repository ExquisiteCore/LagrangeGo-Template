# LagrangeGo-Template

A template for LagrangeGo

[![Go Report Card](https://goreportcard.com/badge/github.com/ExquisiteCore/LagrangeGo-Template)](https://goreportcard.com/report/github.com/ExquisiteCore/LagrangeGo-Template)

基于 [LagrangeGo](https://github.com/LagrangeDev/LagrangeGo) 的Bot 模板参考自[MiraiGo-Template](https://github.com/Logiase/MiraiGo-Template)

## 基础配置

账号配置[application.toml](./application.toml)

```toml
[bot]
# 账号 必填
account = 114514
# 密码 选填
password = "pwd"
```

不配置密码的话将使用扫码登录

## 快速入门

克隆本项目

在[logic/custom_logic.go](./logic/custom_logic.go)注册逻辑

```go
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
```

## 引入的第三方 go module

- [LagrangeGo](https://github.com/LagrangeDev/LagrangeGo)
    核心协议库
- [toml](https://github.com/BurntSushi/toml)
    用于解析配置文件，同时可监听配toml置文件的修改
- [logrus](https://github.com/sirupsen/logrus)
    功能丰富的Logger
