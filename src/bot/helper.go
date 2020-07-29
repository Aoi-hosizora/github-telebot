package bot

import (
	"github.com/Aoi-hosizora/github-telebot/src/logger"
	"gopkg.in/tucnak/telebot.v2"
)

// Handle text endpoint with Message handler.
func (b *bot) handleMessage(endpoint string, handler func(*telebot.Message)) {
	b.Bot.Handle(endpoint, func(m *telebot.Message) {
		logger.Telebot.Receive(endpoint, m)
		handler(m)
	})
}

// Handle inline button endpoint with callback handler.
func (b *bot) handleInline(endpoint *telebot.InlineButton, handler func(*telebot.Callback)) {
	b.Bot.Handle(endpoint, func(c *telebot.Callback) {
		logger.Telebot.Receive(endpoint, c)
		handler(c)
	})
}

// Handle reply button endpoint with callback handler.
func (b *bot) handleReply(endpoint *telebot.ReplyButton, handler func(*telebot.Message)) {
	b.Bot.Handle(endpoint, func(m *telebot.Message) {
		logger.Telebot.Receive(endpoint, m)
		handler(m)
	})
}

// Reply content to a specific message.
func (b *bot) Reply(m *telebot.Message, what interface{}, options ...interface{}) error {
	msg, err := b.Bot.Send(m.Chat, what, options...)
	logger.Telebot.Reply(m, msg, err)
	if err != nil {
		msg, err := b.Bot.Send(m.Chat, "Something went wrong, Please retry.")
		logger.Telebot.Reply(m, msg, err)
	}
	return err
}

// Send content to a specific chat (ByID).
func (b *bot) Send(c *telebot.Chat, what interface{}, options ...interface{}) error {
	msg, err := b.Bot.Send(c, what, options...)
	logger.Telebot.Send(c, msg, err)
	return err
}

// Mirror method from `telebot.Delete`
func (b *bot) Delete(msg telebot.Editable) error {
	return b.Bot.Delete(msg)
}
