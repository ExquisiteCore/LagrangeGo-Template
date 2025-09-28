package bot

import (
	"github.com/LagrangeDev/LagrangeGo/client"
)

// Bot 使用组合模式而不是继承
type Bot struct {
	client        *client.QQClient
	loginMgr      *LoginManager
	authMgr       *AuthManager
	connectionMgr *ConnectionManager
}

// Bot 实例
var QQClient *Bot

// NewBot 创建新的Bot实例
func NewBot(client *client.QQClient) *Bot {
	bot := &Bot{
		client:   client,
		loginMgr: NewLoginManager(client),
		authMgr:  NewAuthManager(client, "sig.bin"),
	}

	// 创建连接管理器并注册默认事件处理器
	bot.connectionMgr = NewConnectionManager(client)
	bot.connectionMgr.RegisterEventHandler(&DefaultConnectionEventHandler{})

	return bot
}

// Client 获取内部客户端
func (b *Bot) Client() *client.QQClient {
	return b.client
}

// Login 登录
func (b *Bot) Login() error {
	return b.loginMgr.Login()
}

// Listen 开始监听连接状态
func (b *Bot) Listen() {
	b.connectionMgr.StartMonitoring()
}

// Stop 停止Bot
func (b *Bot) Stop() {
	b.connectionMgr.StopMonitoring()
}

// Dumpsig 保存签名
func (b *Bot) Dumpsig() {
	b.authMgr.Dumpsig()
}

// GetConnectionManager 获取连接管理器
func (b *Bot) GetConnectionManager() *ConnectionManager {
	return b.connectionMgr
}

// GetLoginManager 获取登录管理器
func (b *Bot) GetLoginManager() *LoginManager {
	return b.loginMgr
}

// GetAuthManager 获取认证管理器
func (b *Bot) GetAuthManager() *AuthManager {
	return b.authMgr
}
