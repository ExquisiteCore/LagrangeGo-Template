package bot

import (
	"bytes"
	"errors"
	"os"
	"time"

	"github.com/ExquisiteCore/LagrangeGo-Template/config"
	"github.com/ExquisiteCore/LagrangeGo-Template/utils"
	"github.com/mdp/qrterminal/v3"
	"github.com/tuotoo/qrcode"

	"github.com/LagrangeDev/LagrangeGo/client"
	"github.com/LagrangeDev/LagrangeGo/client/auth"
	"github.com/LagrangeDev/LagrangeGo/client/packets/wtlogin/qrcodestate"
	"github.com/sirupsen/logrus"
)

// Bot 使用组合模式而不是继承
type Bot struct {
	client   *client.QQClient
	loginMgr *LoginManager
	authMgr  *AuthManager
}

// Bot 实例
var QQClient *Bot

// LoginManager 处理登录相关逻辑
type LoginManager struct {
	client *client.QQClient
}

// AuthManager 处理认证相关逻辑
type AuthManager struct {
	client  *client.QQClient
	sigFile string
}

// NewBot 创建新的Bot实例
func NewBot(client *client.QQClient) *Bot {
	return &Bot{
		client:   client,
		loginMgr: &LoginManager{client: client},
		authMgr:  &AuthManager{client: client, sigFile: "sig.bin"},
	}
}

// Client 获取内部客户端
func (b *Bot) Client() *client.QQClient {
	return b.client
}

// Login 登录
func (b *Bot) Login() error {
	return b.loginMgr.Login()
}

// Dumpsig 保存签名
func (b *Bot) Dumpsig() {
	b.authMgr.Dumpsig()
}

// Listen 监听状态
func (b *Bot) Listen() {
	b.client.DisconnectedEvent.Subscribe(func(client *client.QQClient, event *client.DisconnectedEvent) {
		logrus.Infof("连接已断开：%v", event.Message)
	})
}
func Init(logger *utils.ProtocolLogger) {
	appInfo := auth.AppList["linux"]["3.2.15-30366"]
	//qqClientInstance := client.NewClient(config.GlobalConfig.Bot.Account, appInfo, "https://sign.lagrangecore.org/api/sign/25765")
	qqClientInstance := client.NewClient(config.GlobalConfig.Bot.Account, config.GlobalConfig.Bot.Password)
	qqClientInstance.SetLogger(logger)
	qqClientInstance.UseVersion(appInfo)
	qqClientInstance.AddSignServer(config.GlobalConfig.Bot.SignServer)
	qqClientInstance.UseDevice(auth.NewDeviceInfo(114514))

	QQClient = NewBot(qqClientInstance)
	QQClient.authMgr.loadSig()
}

// Login 登录
func (lm *LoginManager) Login() error {
	//自动登录
	sig := lm.client.Sig()
	if sig != nil {
		err := lm.client.FastLogin()
		if err != nil {
			logrus.Errorln("fast login err:", err)
			logrus.Println("change to qrcode mode.")
		} else {
			return nil
		}
	}

	//获取二维码
	png, _, err := lm.client.FetchQRCodeDefault()
	if err != nil {
		logrus.Errorln("login err:", err)
		return err
	}

	// 打印二维码内容
	qrMatrix, err := qrcode.Decode(bytes.NewReader(png))
	if err != nil {
		logrus.Warnln("二维码内容识别失败，可能无法终端打印:", err)
	} else {
		logrus.Infoln("请使用手机扫码登录：")
		config := qrterminal.Config{
			Level:     qrterminal.M,
			Writer:    os.Stdout,
			BlackChar: qrterminal.WHITE,
			WhiteChar: qrterminal.BLACK,
			QuietZone: 1,
		}
		qrterminal.GenerateWithConfig(qrMatrix.Content, config)
	}

	//轮询登录状态
	var retCode qrcodestate.State
	for {
		retCode, err = lm.client.GetQRCodeResult()
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
	_, err = lm.client.QRCodeLogin()
	if err != nil {
		logrus.Errorln("login err:", err)
		return err
	}
	return nil
}

// loadSig 加载签名文件
func (am *AuthManager) loadSig() {
	data, err := os.ReadFile(am.sigFile)
	if err != nil {
		logrus.Warnln("read sig error:", err)
	} else {
		sig, err := auth.UnmarshalSigInfo(data, true)
		if err != nil {
			logrus.Warnln("load sig error:", err)
		} else {
			am.client.UseSig(sig)
		}
	}
}

// Dumpsig 保存签名
func (am *AuthManager) Dumpsig() {
	data, err := am.client.Sig().Marshal()
	if err != nil {
		logrus.Errorln("marshal sig.bin err:", err)
		return
	}
	err = os.WriteFile(am.sigFile, data, 0644)
	if err != nil {
		logrus.Errorln("write sig.bin err:", err)
		return
	}
	logrus.Infoln("sig saved into sig.bin")
}
