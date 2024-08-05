package config

import (
	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Bot BotConfig
}

// BotConfig 代表TOML文件中的bot部分
type BotConfig struct {
	Account  uint32 `toml:"account"`
	Password string `toml:"password"`
}

// GlobalConfig 默认全局配置
var GlobalConfig *Config

// Init 使用 ./application.toml 初始化全局配置
func Init() {
	GlobalConfig = &Config{}
	_, err := toml.DecodeFile("application.toml", GlobalConfig)
	if err != nil {
		logrus.WithField("config", "GlobalConfig").WithError(err).Panicf("unable to read global config")
	}
}

// InitWithContent 从字节数组中读取配置内容
func InitWithContent(configTOMLContent []byte) {
	_, err := toml.Decode(string(configTOMLContent), GlobalConfig)
	if err != nil {
		panic(err)
	}
}
