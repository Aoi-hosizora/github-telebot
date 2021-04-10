package server

import (
	"github.com/Aoi-hosizora/ahlib-web/xtelebot"
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"github.com/Aoi-hosizora/github-telebot/internal/bot/fsm"
	"github.com/Aoi-hosizora/github-telebot/internal/pkg/config"
	"github.com/Aoi-hosizora/github-telebot/internal/pkg/logger"
	"gopkg.in/tucnak/telebot.v2"
	"strings"
)

// _bot represents the global BotServer.
var _bot *BotServer

func Bot() *BotServer {
	return _bot
}

func SetupBot(b *BotServer) {
	_bot = b
}

type BotServer struct {
	bot  *telebot.Bot
	data *xtelebot.BotData
}

func NewBotServer(bot *telebot.Bot) *BotServer {
	return &BotServer{
		bot:  bot,
		data: xtelebot.NewBotData(xtelebot.WithInitialStatus(fsm.None)),
	}
}

// ===
// Bot
// ===

func (b *BotServer) Start() {
	b.bot.Start()
}

func (b *BotServer) Stop() {
	b.bot.Stop()
}

func (b *BotServer) Edit(msg telebot.Editable, what interface{}, options ...interface{}) (*telebot.Message, error) {
	return b.bot.Edit(msg, what, options...)
}

func (b *BotServer) sendTo(c *telebot.Chat, what interface{}, cb func(*telebot.Message, error), options ...interface{}) error {
	var msg *telebot.Message
	var err error

	for i := 0; i < int(config.Configs().Bot.RetryCount); i++ { // retry
		msg, err = b.bot.Send(c, what, options...)
		cb(msg, err)
		if err == nil {
			return nil
		}

		if flag, ok := what.(string); ok && strings.Contains(err.Error(), "must be escaped") {
			flag = strings.ReplaceAll(flag, "\\", "")
			flag += "\n\nPlease contact to the developer with the message:\n" + err.Error()
			newOptions := make([]interface{}, 1, len(options))
			newOptions[0] = telebot.ModeMarkdown
			for _, opt := range options {
				if opt != telebot.ModeMarkdownV2 {
					newOptions = append(newOptions, opt)
				}
			}
			_, _ = b.bot.Send(c, flag, newOptions...)
			return err
		}
	}

	return err
}

func (b *BotServer) Send(c *telebot.Chat, what interface{}, options ...interface{}) error {
	return b.sendTo(c, what, func(msg *telebot.Message, err error) {
		logger.Send(c, msg, err)
	}, options...)
}

func (b *BotServer) Reply(m *telebot.Message, what interface{}, options ...interface{}) error {
	return b.sendTo(m.Chat, what, func(msg *telebot.Message, err error) {
		logger.Reply(m, msg, err)
	}, options...)
}

func (b *BotServer) SendToChat(chatId int64, what interface{}, options ...interface{}) error {
	chat, err := b.bot.ChatByID(xnumber.I64toa(chatId))
	if err != nil {
		return err
	}

	return b.Send(chat, what, options...)
}

// ======
// Handle
// ======

func (b *BotServer) HandleMessage(endpoint string, handler func(*telebot.Message)) {
	if handler == nil {
		panic("nil handler")
	}
	b.bot.Handle(endpoint, func(m *telebot.Message) {
		logger.Receive(endpoint, m)
		handler(m)
	})
}

func (b *BotServer) HandleInline(endpoint *telebot.InlineButton, handler func(*telebot.Callback)) {
	if handler == nil {
		panic("nil handler")
	}
	b.bot.Handle(endpoint, func(c *telebot.Callback) {
		logger.Receive(endpoint, c.Message)
		handler(c)
	})
}

func (b *BotServer) HandleReply(endpoint *telebot.ReplyButton, handler func(*telebot.Message)) {
	if handler == nil {
		panic("nil handler")
	}
	b.bot.Handle(endpoint, func(m *telebot.Message) {
		logger.Receive(endpoint, m)
		handler(m)
	})
}

// =======
// BotData
// =======

func (b *BotServer) SetStatus(chatID int64, status xtelebot.ChatStatus) {
	b.data.SetStatus(chatID, status)
}

func (b *BotServer) GetStatus(chatID int64) xtelebot.ChatStatus {
	return b.data.GetStatusOrInit(chatID)
}

func (b *BotServer) SetCache(chatID int64, key string, value interface{}) {
	b.data.SetCache(chatID, key, value)
}

func (b *BotServer) GetCache(chatID int64, key string) (interface{}, bool) {
	return b.data.GetCache(chatID, key)
}

func (b *BotServer) RemoveCache(chatID int64, key string) {
	b.data.RemoveCache(chatID, key)
}
