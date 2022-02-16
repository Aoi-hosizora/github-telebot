package controller

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/xtelebot"
	"github.com/Aoi-hosizora/github-telebot/internal/bot/fsm"
	"gopkg.in/tucnak/telebot.v2"
)

const (
	_START           = "Here is github events telebot, developed by @AoiHosizora, send /help for help."
	_NO_ACTION       = "There is no action now."
	_ACTION_CANCELED = "Action \"%s\" has been canceled."
	_UNKNOWN_COMMAND = "Unknown command: %s, send /help for help."
)

// Start /start
func Start(bw *xtelebot.BotWrapper, m *telebot.Message) {
	bw.RespondReply(m, false, _START)
}

// Help /help
func Help(help string) xtelebot.MessageHandler {
	return func(bw *xtelebot.BotWrapper, m *telebot.Message) {
		bw.RespondReply(m, false, help, telebot.ModeMarkdown)
	}
}

// Cancel /cancel
func Cancel(bw *xtelebot.BotWrapper, m *telebot.Message) {
	state := bw.Data().GetStateOr(m.Chat.ID, fsm.None)
	if state == fsm.None {
		bw.RespondReply(m, false, _NO_ACTION, xtelebot.RemoveReplyKeyboard())
	} else {
		bw.Data().SetState(m.Chat.ID, fsm.None)
		s := fmt.Sprintf(_ACTION_CANCELED, fsm.StateString(state))
		bw.RespondReply(m, false, s, xtelebot.RemoveReplyKeyboard())
	}
}

// OnText $on_text
func OnText(bw *xtelebot.BotWrapper, m *telebot.Message) {
	state := bw.Data().GetStateOr(m.Chat.ID, fsm.None)
	handler := bw.Shs().GetHandler(state)
	if handler != nil {
		handler(bw, m)
	} else {
		bw.RespondReply(m, false, fmt.Sprintf(_UNKNOWN_COMMAND, m.Text))
	}
}
