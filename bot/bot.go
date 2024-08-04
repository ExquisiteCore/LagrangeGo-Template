package bot

import (
	"os"

	"github.com/ExquisiteCore/LagrangeGo-Template/config"
	"github.com/ExquisiteCore/LagrangeGo-Template/utils"

	"github.com/LagrangeDev/LagrangeGo/client"
	"github.com/LagrangeDev/LagrangeGo/client/auth"
	"github.com/sirupsen/logrus"
)

// Bot 全局 Bot
type Bot struct {
	*client.QQClient
}

// Bot 实例
var QQClient *Bot

func Init(logger *utils.ProtocolLogger) {
	appInfo := auth.AppList["linux"]["3.2.10-25765"]
	deviceInfo := auth.NewDeviceInfo(114514)
	qqClientInstance := client.NewClient(config.GlobalConfig.Bot.Account, appInfo, "https://sign.lagrangecore.org/api/sign/25765")
	qqClientInstance.SetLogger(logger)
	qqClientInstance.UseDevice(deviceInfo)

	data, err := os.ReadFile("sig.bin")
	if err != nil {
		logrus.Warnln("read sig error:", err)
	} else {
		sig, err := auth.UnmarshalSigInfo(data, true)
		if err != nil {
			logrus.Warnln("load sig error:", err)
		} else {
			qqClientInstance.UseSig(sig)
		}
	}
	QQClient = &Bot{QQClient: qqClientInstance}

}

// Login 登录
func Login() error {
	// 声明 err 变量并进行错误处理
	err := QQClient.Login(config.GlobalConfig.Bot.Password, "qrcode.png")
	if err != nil {
		logrus.Errorln("login err:", err)
		return err
	}
	return nil
}

// 保存sign
func Dumpsig() {
	data, err := QQClient.Sig().Marshal()
	if err != nil {
		logrus.Errorln("marshal sig.bin err:", err)
		return
	}
	err = os.WriteFile("sig.bin", data, 0644)
	if err != nil {
		logrus.Errorln("write sig.bin err:", err)
		return
	}
	logrus.Infoln("sig saved into sig.bin")
}
