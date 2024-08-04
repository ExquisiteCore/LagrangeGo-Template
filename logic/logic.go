package logic

import (
	"github.com/ExquisiteCore/LagrangeGo-Template/bot"
	"github.com/LagrangeDev/LagrangeGo/client"
	"github.com/LagrangeDev/LagrangeGo/message"
)

// Bot逻辑
func Startlogic() {
	bot.QQClient.PrivateMessageEvent.Subscribe(func(client *client.QQClient, event *message.PrivateMessage) {
		_, err := client.SendPrivateMessage(event.Sender.Uin, []message.IMessageElement{message.NewText("你好")})
		if err != nil {
			return
		}
	})
}
