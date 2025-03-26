package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/ExquisiteCore/LagrangeGo-Template/bot"
	"github.com/ExquisiteCore/LagrangeGo-Template/config"
	"github.com/ExquisiteCore/LagrangeGo-Template/logic"
	"github.com/ExquisiteCore/LagrangeGo-Template/utils"
)

// 创建 protocolLogger 实例
var logger = utils.ProtocolLogger{}

func init() {
	config.Init()
	utils.Init()
	bot.Init(&logger)
}

func main() {

	bot.Login()

	bot.Listen()

	logic.RegisterCustomLogic()

	logic.SetupLogic()

	defer bot.QQClient.Release()

	defer bot.Dumpsig()

	// setup the main stop channel
	mc := make(chan os.Signal, 2)
	signal.Notify(mc, os.Interrupt, syscall.SIGTERM)
	for {
		switch <-mc {
		case os.Interrupt, syscall.SIGTERM:
			return
		}
	}
}
