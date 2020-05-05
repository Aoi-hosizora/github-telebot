package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var (
	Configs *Config
)

type Config struct {
	Mode           string          `yaml:"mode"`
	TelegramConfig *TelegramConfig `yaml:"telegram"`
	TaskConfig     *TaskConfig     `yaml:"task"`
	MysqlConfig    *MysqlConfig    `yaml:"mysql"`
}

type TelegramConfig struct {
	PollerTimeout int    `yaml:"poller-timeout"`
	BotToken      string `yaml:"bot-token"`
}

type TaskConfig struct {
	PollingDuration int `yaml:"polling-duration"`
}

type MysqlConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
	LogMode  bool   `json:"log-mode"`
}

func LoadConfig(configPath string) error {
	f, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}
	config := &Config{}
	err = yaml.Unmarshal(f, config)
	if err != nil {
		return err
	}

	Configs = config
	return nil
}
