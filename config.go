package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type ServerConfig struct {
	PollingDuration uint32 `yaml:"polling-duration"`
}

type TelegramConfig struct {
	BotToken  string `yaml:"bot-token"`
	ChannelId string `yaml:"channel-id"`
}

type GithubConfig struct {
	Token string `yaml:"token"`
}

type Config struct {
	ServerConfig   *ServerConfig   `yaml:"server"`
	TelegramConfig *TelegramConfig `yaml:"telegram"`
	GithubConfig   *GithubConfig   `yaml:"github"`
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
