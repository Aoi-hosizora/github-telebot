package logger

import (
	"github.com/Aoi-hosizora/ahlib-more/xlogrus"
	"github.com/Aoi-hosizora/ahlib-web/xtelebot"
	"github.com/Aoi-hosizora/github-telebot/internal/pkg/config"
	"github.com/sirupsen/logrus"
	"gopkg.in/tucnak/telebot.v2"
	"time"
)

// _logger represents the global logrus.Logger.
var _logger *logrus.Logger

func Logger() *logrus.Logger {
	return _logger
}

func Setup() error {
	logger := logrus.New()
	logLevel := logrus.WarnLevel
	if config.Configs().Meta.RunMode == "debug" {
		logLevel = logrus.DebugLevel
	}

	logger.SetLevel(logLevel)
	logger.SetReportCaller(false)
	logger.SetFormatter(&xlogrus.SimpleFormatter{TimestampFormat: time.RFC3339})
	logger.AddHook(xlogrus.NewRotateLogHook(&xlogrus.RotateLogConfig{
		Filename:         config.Configs().Meta.LogName,
		FilenameTimePart: ".%Y%m%d.log",
		LinkFileName:     config.Configs().Meta.LogName + ".log",
		Level:            logLevel,
		Formatter:        &logrus.JSONFormatter{TimestampFormat: time.RFC3339},
		MaxAge:           15 * 24 * time.Hour,
		RotationTime:     24 * time.Hour,
	}))

	_logger = logger
	return nil
}

// Receive is a global function to log receive message to logrus.Logger.
func Receive(endpoint interface{}, message *telebot.Message) {
	xtelebot.LogReceiveToLogrus(_logger, endpoint, message)
}

// Reply is a global function to log reply message to logrus.Logger.
func Reply(received, replied *telebot.Message, err error) {
	xtelebot.LogReplyToLogrus(_logger, received, replied, err)
}

// Send is a global function to log send message to logrus.Logger.
func Send(chat *telebot.Chat, sent *telebot.Message, err error) {
	xtelebot.LogSendToLogrus(_logger, chat, sent, err)
}
