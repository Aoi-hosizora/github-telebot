package main

import (
	"flag"
	"github.com/Aoi-hosizora/ah-tgbot/src/config"
	"github.com/Aoi-hosizora/ah-tgbot/src/server"
	"log"
)

var (
	help       bool
	configPath string
)

func init() {
	flag.BoolVar(&help, "h", false, "show help")
	flag.StringVar(&configPath, "config", "./src/config/config.yaml", "change the config path")
}

func main() {
	flag.Parse()
	if help {
		flag.Usage()
	} else {
		run()
	}
}

func run() {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalln("Failed to load config file:", err)
	}
	bs := server.NewBotServer(cfg)
	defer func() {
		bs.Bot.Stop()
	}()
	bs.Serve()
}
