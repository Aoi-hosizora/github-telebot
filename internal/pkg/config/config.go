package config

import (
	"github.com/Aoi-hosizora/ahlib-web/xvalidator"
	"github.com/Aoi-hosizora/ahlib/xreflect"
	"gopkg.in/yaml.v2"
	"os"
)

// _configs represents the global config.Config.
var _configs *Config

func Configs() *Config {
	return _configs
}

type Config struct {
	Meta   *MetaConfig   `yaml:"meta"   validate:"required"`
	Task   *TaskConfig   `yaml:"task"   validate:"required"`
	SQLite *SQLiteConfig `yaml:"sqlite" validate:"required"`
	Redis  *RedisConfig  `yaml:"redis"  validate:"required"`
}

type MetaConfig struct {
	Token   string `yaml:"token"    validate:"required"`
	RunMode string `yaml:"run-mode" default:"debug"`
	LogName string `yaml:"log-name" default:"./logs/console"`

	PollerTimeout uint64 `yaml:"poller-timeout" default:"5" validate:"gt=0"`
}

type TaskConfig struct {
	ActivityCron string `yaml:"activity-cron" validate:"required"`
	IssueCron    string `yaml:"issue-cron"    validate:"required"`
}

type SQLiteConfig struct {
	Database string `yaml:"database" validate:"required"`
	LogMode  bool   `yaml:"log-mode"`
}

type RedisConfig struct {
	Host     string `yaml:"host" default:"127.0.0.1"`
	Port     int32  `yaml:"port" default:"6379"`
	DB       int32  `yaml:"db"   validate:"required"`
	Password string `yaml:"password"`
	LogMode  bool   `yaml:"log-mode"`

	DialTimeout  *int32 `yaml:"dial-timeout"  validate:"omitempty,gt=0"`
	ReadTimeout  *int32 `yaml:"read-timeout"  validate:"omitempty,gt=0"`
	WriteTimeout *int32 `yaml:"write-timeout" validate:"omitempty,gt=0"`
	MaxOpens     *int32 `yaml:"max-opens"     validate:"omitempty,gt=0"`
	MinIdles     *int32 `yaml:"min-idles"     validate:"omitempty,gte=0"`
	MaxLifetime  *int32 `yaml:"max-lifetime"  validate:"omitempty,gt=0"`
	MaxIdletime  *int32 `yaml:"max-idletime"  validate:"omitempty,gt=0"`
}

var _debugMode = true

func IsDebugMode() bool {
	return _debugMode
}

func Load(path string) error {
	f, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	cfg := &Config{}
	if err = yaml.Unmarshal(f, cfg); err != nil {
		return err
	}
	if _, err = xreflect.FillDefaultFields(cfg); err != nil {
		return err
	}
	if err = validateConfig(cfg); err != nil {
		return err
	}

	_debugMode = cfg.Meta.RunMode == "debug"
	_configs = cfg
	return nil
}

func validateConfig(cfg *Config) error {
	val := xvalidator.NewMessagedValidator()
	val.SetValidateTagName("validate")
	val.SetMessageTagName("message")
	val.UseTagAsFieldName("yaml", "json")
	if err := val.ValidateStruct(cfg); err != nil {
		ut, _ := xvalidator.ApplyTranslator(val.ValidateEngine(), xvalidator.EnLocaleTranslator(), xvalidator.EnTranslationRegisterFunc())
		return xvalidator.MapToError(err.(*xvalidator.MultiFieldsError).Translate(ut, false))
	}
	return nil
}
