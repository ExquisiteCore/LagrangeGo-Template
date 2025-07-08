package bot

import (
	"context"
	"errors"
	"time"

	"github.com/LagrangeDev/LagrangeGo/client"
	"github.com/sirupsen/logrus"
)

// LoginStrategy 登录策略接口
type LoginStrategy interface {
	Login(ctx context.Context, client *client.QQClient) error
	GetStrategyName() string
}

// LoginContext 登录上下文
type LoginContext struct {
	MaxRetries int
	RetryDelay time.Duration
	Timeout    time.Duration
}

// DefaultLoginContext 默认登录上下文
func DefaultLoginContext() *LoginContext {
	return &LoginContext{
		MaxRetries: 3,
		RetryDelay: 3 * time.Second,
		Timeout:    5 * time.Minute,
	}
}

// LoginManager 登录管理器
type LoginManager struct {
	client     *client.QQClient
	strategies []LoginStrategy
	context    *LoginContext
}

// NewLoginManager 创建新的登录管理器
func NewLoginManager(client *client.QQClient) *LoginManager {
	lm := &LoginManager{
		client:  client,
		context: DefaultLoginContext(),
	}

	// 注册默认登录策略
	lm.RegisterStrategy(&FastLoginStrategy{})
	lm.RegisterStrategy(&QRCodeLoginStrategy{})

	return lm
}

// RegisterStrategy 注册登录策略
func (lm *LoginManager) RegisterStrategy(strategy LoginStrategy) {
	lm.strategies = append(lm.strategies, strategy)
}

// Login 执行登录
func (lm *LoginManager) Login() error {
	ctx, cancel := context.WithTimeout(context.Background(), lm.context.Timeout)
	defer cancel()

	for _, strategy := range lm.strategies {
		logrus.Infof("尝试使用 %s 登录", strategy.GetStrategyName())

		err := lm.tryLoginWithRetry(ctx, strategy)
		if err == nil {
			logrus.Infof("使用 %s 登录成功", strategy.GetStrategyName())
			return nil
		}

		logrus.Warnf("使用 %s 登录失败: %v", strategy.GetStrategyName(), err)
	}

	return errors.New("所有登录策略都失败了")
}

// tryLoginWithRetry 带重试的登录尝试
func (lm *LoginManager) tryLoginWithRetry(ctx context.Context, strategy LoginStrategy) error {
	var lastErr error

	for i := 0; i < lm.context.MaxRetries; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		err := strategy.Login(ctx, lm.client)
		if err == nil {
			return nil
		}

		lastErr = err
		if i < lm.context.MaxRetries-1 {
			logrus.Warnf("第 %d 次尝试失败: %v，%v 后重试", i+1, err, lm.context.RetryDelay)
			time.Sleep(lm.context.RetryDelay)
		}
	}

	return lastErr
}

// FastLoginStrategy 快速登录策略
type FastLoginStrategy struct{}

func (s *FastLoginStrategy) GetStrategyName() string {
	return "快速登录"
}

func (s *FastLoginStrategy) Login(ctx context.Context, client *client.QQClient) error {
	sig := client.Sig()
	if sig == nil {
		return errors.New("没有可用的签名信息")
	}

	return client.FastLogin()
}

// QRCodeLoginStrategy 二维码登录策略
type QRCodeLoginStrategy struct {
	qrProcessor *QRCodeProcessor
}

func (s *QRCodeLoginStrategy) GetStrategyName() string {
	return "二维码登录"
}

func (s *QRCodeLoginStrategy) Login(ctx context.Context, client *client.QQClient) error {
	if s.qrProcessor == nil {
		s.qrProcessor = NewQRCodeProcessor()
	}

	// 获取二维码
	png, _, err := client.FetchQRCodeDefault()
	if err != nil {
		return err
	}

	// 显示二维码
	err = s.qrProcessor.DisplayQRCode(png)
	if err != nil {
		logrus.Warnf("二维码显示失败: %v", err)
	}

	// 轮询登录状态
	return s.pollLoginStatus(ctx, client)
}

// pollLoginStatus 轮询登录状态
func (s *QRCodeLoginStrategy) pollLoginStatus(ctx context.Context, client *client.QQClient) error {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			retCode, err := client.GetQRCodeResult()
			if err != nil {
				return err
			}

			if !retCode.Waitable() {
				if !retCode.Success() {
					return errors.New(retCode.Name())
				}

				// 执行二维码登录
				_, err = client.QRCodeLogin()
				return err
			}
		}
	}
}
