package bot

import (
	"bytes"
	"os"
	"io"

	"github.com/mdp/qrterminal/v3"
	"github.com/tuotoo/qrcode"
	"github.com/sirupsen/logrus"
	"rsc.io/qr"
)

// QRCodeProcessor 二维码处理器
type QRCodeProcessor struct {
	config *QRCodeConfig
}

// QRCodeConfig 二维码配置
type QRCodeConfig struct {
	Level     qr.Level
	Writer    io.Writer
	BlackChar string
	WhiteChar string
	QuietZone int
}

// DefaultQRCodeConfig 默认二维码配置
func DefaultQRCodeConfig() *QRCodeConfig {
	return &QRCodeConfig{
		Level:     qr.M,
		Writer:    os.Stdout,
		BlackChar: qrterminal.WHITE,
		WhiteChar: qrterminal.BLACK,
		QuietZone: 1,
	}
}

// NewQRCodeProcessor 创建新的二维码处理器
func NewQRCodeProcessor() *QRCodeProcessor {
	return &QRCodeProcessor{
		config: DefaultQRCodeConfig(),
	}
}

// NewQRCodeProcessorWithConfig 使用自定义配置创建二维码处理器
func NewQRCodeProcessorWithConfig(config *QRCodeConfig) *QRCodeProcessor {
	return &QRCodeProcessor{
		config: config,
	}
}

// DisplayQRCode 显示二维码
func (qp *QRCodeProcessor) DisplayQRCode(pngData []byte) error {
	// 解码二维码内容
	qrMatrix, err := qrcode.Decode(bytes.NewReader(pngData))
	if err != nil {
		return err
	}
	
	logrus.Infoln("请使用手机扫码登录：")
	
	// 生成终端二维码
	config := qrterminal.Config{
		Level:     qp.config.Level,
		Writer:    qp.config.Writer,
		BlackChar: qp.config.BlackChar,
		WhiteChar: qp.config.WhiteChar,
		QuietZone: qp.config.QuietZone,
	}
	
	qrterminal.GenerateWithConfig(qrMatrix.Content, config)
	return nil
}

// SaveQRCodeToFile 保存二维码到文件
func (qp *QRCodeProcessor) SaveQRCodeToFile(pngData []byte, filename string) error {
	return os.WriteFile(filename, pngData, 0644)
}

// GetQRCodeContent 获取二维码内容
func (qp *QRCodeProcessor) GetQRCodeContent(pngData []byte) (string, error) {
	qrMatrix, err := qrcode.Decode(bytes.NewReader(pngData))
	if err != nil {
		return "", err
	}
	return qrMatrix.Content, nil
}