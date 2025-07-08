package bot

import (
	"os"

	"github.com/LagrangeDev/LagrangeGo/client"
	"github.com/LagrangeDev/LagrangeGo/client/auth"
	"github.com/sirupsen/logrus"
)

// AuthManager 处理认证相关逻辑
type AuthManager struct {
	client  *client.QQClient
	sigFile string
}

// NewAuthManager 创建新的认证管理器
func NewAuthManager(client *client.QQClient, sigFile string) *AuthManager {
	return &AuthManager{
		client:  client,
		sigFile: sigFile,
	}
}

// LoadSig 加载签名文件
func (am *AuthManager) LoadSig() {
	data, err := os.ReadFile(am.sigFile)
	if err != nil {
		logrus.Warnln("read sig error:", err)
		return
	}
	
	sig, err := auth.UnmarshalSigInfo(data, true)
	if err != nil {
		logrus.Warnln("load sig error:", err)
		return
	}
	
	am.client.UseSig(sig)
	logrus.Infoln("签名文件加载成功")
}

// Dumpsig 保存签名
func (am *AuthManager) Dumpsig() {
	if am.client.Sig() == nil {
		logrus.Warnln("没有可用的签名信息")
		return
	}
	
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

// HasValidSig 检查是否有有效的签名
func (am *AuthManager) HasValidSig() bool {
	return am.client.Sig() != nil
}

// GetSigFile 获取签名文件路径
func (am *AuthManager) GetSigFile() string {
	return am.sigFile
}

// SetSigFile 设置签名文件路径
func (am *AuthManager) SetSigFile(sigFile string) {
	am.sigFile = sigFile
}