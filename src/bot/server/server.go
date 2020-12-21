package server

import (
	"github.com/Aoi-hosizora/ahlib-web/xtelebot"
	"github.com/Aoi-hosizora/github-telebot/src/bot/fsm"
	"github.com/Aoi-hosizora/github-telebot/src/logger"
	"gopkg.in/tucnak/telebot.v2"
)

var Bot *BotServer

type BotServer struct {
	bot       *telebot.Bot
	UsersData *xtelebot.UsersData
}

func NewBotServer(bot *telebot.Bot) *BotServer {
	return &BotServer{
		bot:       bot,
		UsersData: xtelebot.NewUsersData(fsm.None),
	}
}

func (b *BotServer) Start() {
	b.bot.Start()
}

func (b *BotServer) Stop() {
	b.bot.Stop()
}

func (b *BotServer) Delete(msg telebot.Editable) error {
	return b.bot.Delete(msg)
}

func (b *BotServer) ChatByID(id string) (*telebot.Chat, error) {
	return b.bot.ChatByID(id)
}

// Handle string endpoint with telebot.Message handler.
func (b *BotServer) HandleMessage(endpoint string, handler func(*telebot.Message)) {
	b.bot.Handle(endpoint, func(m *telebot.Message) {
		logger.Telebot.Receive(endpoint, m)
		handler(m)
	})
}

// Handle telebot.InlineButton endpoint with telebot.Callback handler.
func (b *BotServer) HandleInline(endpoint *telebot.InlineButton, handler func(*telebot.Callback)) {
	b.bot.Handle(endpoint, func(c *telebot.Callback) {
		logger.Telebot.Receive(endpoint, c)
		handler(c)
	})
}

// Handle telebot.ReplyButton endpoint with telebot.Message handler.
func (b *BotServer) HandleReply(endpoint *telebot.ReplyButton, handler func(*telebot.Message)) {
	b.bot.Handle(endpoint, func(m *telebot.Message) {
		logger.Telebot.Receive(endpoint, m)
		handler(m)
	})
}

// Reply content to a specific message.
func (b *BotServer) Reply(m *telebot.Message, what interface{}, options ...interface{}) error {
	msg, err := b.bot.Send(m.Chat, what, options...)
	logger.Telebot.Reply(m, msg, err)
	if err != nil {
		msg, err := b.bot.Send(m.Chat, "Something went wrong, Please retry.")
		logger.Telebot.Reply(m, msg, err)
	}
	return err
}

// Send content to a specific chat (ByID).
func (b *BotServer) Send(c *telebot.Chat, what interface{}, options ...interface{}) error {
	msg, err := b.bot.Send(c, what, options...)
	logger.Telebot.Send(c, msg, err)
	return err
}
