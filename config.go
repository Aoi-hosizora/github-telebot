package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	ServerConfig *struct {
		PollingDuration uint32 `yaml:"polling-duration"`
	} `yaml:"server"`

	TelegramConfig *struct {
		BotToken  string `yaml:"bot-token"`
		ChannelId string `yaml:"channel-id"`
	} `yaml:"telegram"`

	GithubConfig *struct {
		Username string `yaml:"username"`
		Token    string `yaml:"token"`
	} `yaml:"github"`
}

func loadConfig(configPath string) (*Config, error) {
	f, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	config := new(Config)
	err = yaml.Unmarshal(f, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
