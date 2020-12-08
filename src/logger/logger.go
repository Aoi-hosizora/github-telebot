package logger

import (
	"github.com/Aoi-hosizora/ahlib-more/xlogrus"
	"github.com/Aoi-hosizora/ahlib-web/xtelebot"
	"github.com/Aoi-hosizora/github-telebot/src/config"
	"github.com/sirupsen/logrus"
	"time"
)

var (
	Logger  *logrus.Logger
	Telebot *xtelebot.TelebotLogrus
)

func Setup() error {
	Logger = logrus.New()
	logLevel := logrus.WarnLevel
	if config.Configs.Meta.RunMode == "debug" {
		logLevel = logrus.DebugLevel
	}

	Logger.SetLevel(logLevel)
	Logger.SetReportCaller(false)
	Logger.SetFormatter(&xlogrus.CustomFormatter{TimestampFormat: time.RFC3339})
	Logger.AddHook(xlogrus.NewRotateLogHook(&xlogrus.RotateLogConfig{
		MaxAge:       15 * 24 * time.Hour,
		RotationTime: 24 * time.Hour,
		Filepath:     config.Configs.Meta.LogPath,
		Filename:     config.Configs.Meta.LogName,
		Level:        logLevel,
		Formatter:    &logrus.JSONFormatter{TimestampFormat: time.RFC3339},
	}))

	Telebot = xtelebot.NewTelebotLogrus(Logger, true)

	return nil
}
