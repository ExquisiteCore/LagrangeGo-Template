package bot

import (
	"os"

	"github.com/ExquisiteCore/LagrangeGo-Template/utils"
	"github.com/LagrangeDev/LagrangeGo/client"
	"github.com/LagrangeDev/LagrangeGo/client/auth"
)

// AuthManager 处理认证相关逻辑
type AuthManager struct {
	client  *client.QQClient
	sigFile string
	logger  utils.Logger
}

// NewAuthManager 创建新的认证管理器
func NewAuthManager(client *client.QQClient, sigFile string) *AuthManager {
	return &AuthManager{
		client:  client,
		sigFile: sigFile,
		logger:  utils.GetLogger().WithField("module", "auth"),
	}
}

// LoadSig 加载签名文件
func (am *AuthManager) LoadSig() {
	data, err := os.ReadFile(am.sigFile)
	if err != nil {
		am.logger.Warnf("读取签名文件失败: %v", err)
		return
	}
	
	sig, err := auth.UnmarshalSigInfo(data, true)
	if err != nil {
		am.logger.Warnf("解析签名文件失败: %v", err)
		return
	}
	
	am.client.UseSig(sig)
	am.logger.Info("签名文件加载成功")
}

// Dumpsig 保存签名
func (am *AuthManager) Dumpsig() {
	if am.client.Sig() == nil {
		am.logger.Warn("没有可用的签名信息")
		return
	}
	
	data, err := am.client.Sig().Marshal()
	if err != nil {
		am.logger.Errorf("序列化签名失败: %v", err)
		return
	}
	
	err = os.WriteFile(am.sigFile, data, 0644)
	if err != nil {
		am.logger.Errorf("写入签名文件失败: %v", err)
		return
	}
	
	am.logger.Info("签名文件保存成功")
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