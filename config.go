package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	TelegramConfig *struct{ Token string `yaml:"token"` } `yaml:"telegram"`
	GithubConfig   *struct{ Token string `yaml:"token"` } `yaml:"github"`
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
