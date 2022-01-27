package controller

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/xtelebot"
	"github.com/Aoi-hosizora/github-telebot/internal/bot/fsm"
	"gopkg.in/tucnak/telebot.v2"
)

const (
	_START = "Here is github events telebot, developed by @AoiHosizora, send /help for help."

	_HELP = `*Start*
/start - show start message
/help - show this help message
/cancel - cancel the last action

*Subscribe*
/subscribe - subscribe with a new github account
/unsubscribe - unsubscribe the current github account
/me - show the subscribed user's information

*Option*
/allowissue - allow bot to send issue events
/disallowissue - disallow bot to send issue events
/enablesilent - send message with no notification
/disablesilent - send message with notification
/enablepreview - enable preview for links
/disablepreview - disable preview for links

*Event*
/activity - show the first page of activity events
/activity N - show the N-th page of activity events
/issue - show the first page of issue events
/issue N - show the N-th page of issue events

*Bug report*
https://github.com/Aoi-hosizora/github-telebot/issues`

	_NO_ACTION       = "There is no action now."
	_ACTION_CANCELED = "Action \"%s\" has been canceled."

	_UNKNOWN_COMMAND = "Unknown command: %s, send /help for help."
)

// Start /start
func Start(bw *xtelebot.BotWrapper, m *telebot.Message) {
	bw.ReplyTo(m, _START)
}

// Help /help
func Help(bw *xtelebot.BotWrapper, m *telebot.Message) {
	bw.ReplyTo(m, _HELP, telebot.ModeMarkdown)
}

// Cancel /cancel
func Cancel(bw *xtelebot.BotWrapper, m *telebot.Message) {
	markup := &telebot.ReplyMarkup{ReplyKeyboardRemove: true}
	if state := bw.Data().GetStateOr(m.Chat.ID, fsm.None); state == fsm.None {
		bw.ReplyTo(m, _NO_ACTION, markup)
	} else {
		bw.Data().SetState(m.Chat.ID, fsm.None)
		bw.ReplyTo(m, fmt.Sprintf(_ACTION_CANCELED, fsm.StateString(state)), markup)
	}
}

// OnText $on_text
func OnText(bw *xtelebot.BotWrapper, m *telebot.Message) {
	state := bw.Data().GetStateOr(m.Chat.ID, fsm.None)
	handler := fsm.GetStateHandler(state)
	if handler != nil {
		handler(bw, m)
	} else {
		bw.ReplyTo(m, fmt.Sprintf(_UNKNOWN_COMMAND, m.Text))
	}
}
