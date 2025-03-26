package bot

import (
	"errors"
	"os"
	"time"

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
	appInfo := auth.AppList["linux"]["3.2.15-30366"]
	//qqClientInstance := client.NewClient(config.GlobalConfig.Bot.Account, appInfo, "https://sign.lagrangecore.org/api/sign/25765")
	qqClientInstance := client.NewClient(config.GlobalConfig.Bot.Account, config.GlobalConfig.Bot.Password)
	qqClientInstance.SetLogger(logger)
	qqClientInstance.UseVersion(appInfo)
	qqClientInstance.AddSignServer(config.GlobalConfig.Bot.SignServer)
	qqClientInstance.UseDevice(auth.NewDeviceInfo(114514))

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
	//自动登录
	sig := QQClient.Sig()
	if sig != nil {
		err := QQClient.FastLogin()
		if err != nil {
			logrus.Errorln("fast login err:", err)
			logrus.Println("change to qrcode mode.")
		} else {
			return nil
		}
	}

	//获取二维码
	png, _, err := QQClient.FetchQRCodeDefault()
	if err != nil {
		logrus.Errorln("login err:", err)
		return err
	}

	//保存本地二维码
	qrcodePath := "qrcode.png"
	err = os.WriteFile(qrcodePath, png, 0644)
	if err != nil {
		logrus.Errorln("write qrcode err:", err)
		return err
	}
	//打印二维码
	logrus.Infof("qrcode saved to %s", qrcodePath)
	//轮询登录状态
	for {
		retCode, err := QQClient.GetQRCodeResult()
		if err != nil {
			logrus.Errorln(err)
			return err
		}
		// 等待扫码
		if retCode.Waitable() {
			time.Sleep(3 * time.Second)
			continue
		}
		if !retCode.Success() {
			return errors.New(retCode.Name())
		}
		break
	}
	_, err = QQClient.QRCodeLogin()
	if err != nil {
		logrus.Errorln("login err:", err)
		return err
	}
	return nil
}

// 监听状态
func Listen() {
	QQClient.DisconnectedEvent.Subscribe(func(client *client.QQClient, event *client.DisconnectedEvent) {
		logrus.Infof("连接已断开：%v", event.Message)
	})
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
