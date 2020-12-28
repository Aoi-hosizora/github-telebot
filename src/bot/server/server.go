package server

import (
	"github.com/Aoi-hosizora/ahlib-web/xtelebot"
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"github.com/Aoi-hosizora/github-telebot/src/bot/fsm"
	"github.com/Aoi-hosizora/github-telebot/src/logger"
	"gopkg.in/tucnak/telebot.v2"
)

// Bot is a global variable.
var Bot *BotServer

type BotServer struct {
	bot  *telebot.Bot
	data *xtelebot.UsersData
}

func NewBotServer(bot *telebot.Bot) *BotServer {
	return &BotServer{
		bot:  bot,
		data: xtelebot.NewUsersData(fsm.None),
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

func (b *BotServer) Delete(msg telebot.Editable) error {
	return b.bot.Delete(msg)
}

func (b *BotServer) Reply(m *telebot.Message, what interface{}, options ...interface{}) error {
	var msg *telebot.Message
	var err error
	for i := 0; i < 5; i++ { // retry
		msg, err = b.bot.Send(m.Chat, what, options...)
		logger.Telebot.Reply(m, msg, err)
		if err == nil {
			break
		}
	}
	return err
}

func (b *BotServer) Send(c *telebot.Chat, what interface{}, options ...interface{}) error {
	msg, err := b.bot.Send(c, what, options...)
	logger.Telebot.Send(c, msg, err)
	return err
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

// Handle string endpoint with telebot.Message handler.
func (b *BotServer) HandleMessage(endpoint string, handler func(*telebot.Message)) {
	if handler == nil {
		panic("nil handler")
	}
	b.bot.Handle(endpoint, func(m *telebot.Message) {
		logger.Telebot.Receive(endpoint, m)
		handler(m)
	})
}

// Handle telebot.InlineButton endpoint with telebot.Callback handler.
func (b *BotServer) HandleInline(endpoint *telebot.InlineButton, handler func(*telebot.Callback)) {
	if handler == nil {
		panic("nil handler")
	}
	b.bot.Handle(endpoint, func(c *telebot.Callback) {
		logger.Telebot.Receive(endpoint, c)
		handler(c)
	})
}

// Handle telebot.ReplyButton endpoint with telebot.Message handler.
func (b *BotServer) HandleReply(endpoint *telebot.ReplyButton, handler func(*telebot.Message)) {
	if handler == nil {
		panic("nil handler")
	}
	b.bot.Handle(endpoint, func(m *telebot.Message) {
		logger.Telebot.Receive(endpoint, m)
		handler(m)
	})
}

// =========
// UsersData
// =========

func (b *BotServer) SetStatus(chatID int64, status xtelebot.UserStatus) {
	b.data.SetStatus(chatID, status)
}

func (b *BotServer) GetStatus(chatID int64) xtelebot.UserStatus {
	return b.data.GetStatus(chatID)
}

func (b *BotServer) SetCache(chatID int64, key string, value interface{}) {
	b.data.SetCache(chatID, key, value)
}

func (b *BotServer) GetCache(chatID int64, key string) interface{} {
	return b.data.GetCache(chatID, key)
}

func (b *BotServer) DeleteCache(chatID int64, key string) {
	b.data.DeleteCache(chatID, key)
}
