package server

import (
	"github.com/Aoi-hosizora/ahlib-web/xtelebot"
	"github.com/Aoi-hosizora/github-telebot/src/bot/fsm"
	"github.com/Aoi-hosizora/github-telebot/src/logger"
	"gopkg.in/tucnak/telebot.v2"
)

var Bot *BotServer

type BotServer struct {
	Bot           *telebot.Bot
	UsersData     *xtelebot.UsersData
	InlineButtons map[string]*telebot.InlineButton
	ReplyButtons  map[string]*telebot.ReplyButton
}

func NewBotServer(bot *telebot.Bot) *BotServer {
	return &BotServer{
		Bot:           bot,
		UsersData:     xtelebot.NewUsersData(fsm.None),
		InlineButtons: make(map[string]*telebot.InlineButton),
		ReplyButtons:  make(map[string]*telebot.ReplyButton),
	}
}

func (b *BotServer) Start() {
	b.Bot.Start()
}

func (b *BotServer) Stop() {
	b.Bot.Stop()
}

func (b *BotServer) Delete(msg telebot.Editable) error {
	return b.Bot.Delete(msg)
}

// Handle text endpoint with Message handler.
func (b *BotServer) HandleMessage(endpoint string, handler func(*telebot.Message)) {
	b.Bot.Handle(endpoint, func(m *telebot.Message) {
		logger.Telebot.Receive(endpoint, m)
		handler(m)
	})
}

// Handle inline button endpoint with callback handler.
func (b *BotServer) HandleInline(endpoint *telebot.InlineButton, handler func(*telebot.Callback)) {
	b.Bot.Handle(endpoint, func(c *telebot.Callback) {
		logger.Telebot.Receive(endpoint, c)
		handler(c)
	})
}

// Handle reply button endpoint with callback handler.
func (b *BotServer) HandleReply(endpoint *telebot.ReplyButton, handler func(*telebot.Message)) {
	b.Bot.Handle(endpoint, func(m *telebot.Message) {
		logger.Telebot.Receive(endpoint, m)
		handler(m)
	})
}

// Reply content to a specific message.
func (b *BotServer) Reply(m *telebot.Message, what interface{}, options ...interface{}) error {
	msg, err := b.Bot.Send(m.Chat, what, options...)
	logger.Telebot.Reply(m, msg, err)
	if err != nil {
		msg, err := b.Bot.Send(m.Chat, "Something went wrong, Please retry.")
		logger.Telebot.Reply(m, msg, err)
	}
	return err
}

// Send content to a specific chat (ByID).
func (b *BotServer) Send(c *telebot.Chat, what interface{}, options ...interface{}) error {
	msg, err := b.Bot.Send(c, what, options...)
	logger.Telebot.Send(c, msg, err)
	return err
}
