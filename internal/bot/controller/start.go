package controller

import (
	"fmt"
	"github.com/Aoi-hosizora/github-telebot/internal/bot/fsm"
	"github.com/Aoi-hosizora/github-telebot/internal/bot/server"
	"gopkg.in/tucnak/telebot.v2"
)

const (
	START = "Here is AoiHosizora's github telebot, developed by @AoiHosizora. Use /help to show help message."
	HELP  = `*Commands*
/start - show start message
/help - show this help message
/cancel - cancel the last action

*Account*
/bind - bind with a new github account
/unbind - unbind an old github account
/me - show the bind user's information
/enablesilent - enable bot silence send
/disablesilent - disable bot silence send

*Filter*
/allowissue - allow bot to send issue events
/disallowissue - disallow bot to send issue events
/listfilter - list all notify filters
/addfilter - add a notify filter
/deletefilter - delete a notify filter

*Events*
/activity - show the first page of activity events
/activitypage - show the nth page of activity events
/issue - show the first page of issue events
/issuepage - show the nth page of issue events

*Bug report*
https://github.com/Aoi-hosizora/github-telebot/issues/new`

	NO_ACTION       = "There is no action now."
	ACTION_CANCELED = "Current action has been canceled."
	UNKNOWN_COMMAND = "Unknown command: %s. Send /help to see help."
)

// /start
func StartCtrl(m *telebot.Message) {
	_ = server.Bot().Reply(m, START)
}

// /help
func HelpCtrl(m *telebot.Message) {
	_ = server.Bot().Reply(m, HELP, telebot.ModeMarkdown)
}

// /cancel
func CancelCtrl(m *telebot.Message) {
	if server.Bot().GetStatus(m.Chat.ID) == fsm.None {
		_ = server.Bot().Reply(m, NO_ACTION, &telebot.ReplyMarkup{
			ReplyKeyboardRemove: true,
		})
	} else {
		server.Bot().SetStatus(m.Chat.ID, fsm.None)
		_ = server.Bot().Reply(m, ACTION_CANCELED, &telebot.ReplyMarkup{
			ReplyKeyboardRemove: true,
		})
	}
}

// button.InlineBtnCancel
func InlineBtnCancelCtrl(c *telebot.Callback) {
	m := c.Message
	_, _ = server.Bot().Edit(m, fmt.Sprintf("%s (canceled)", m.Text))
}

// $on_text
func OnTextCtrl(m *telebot.Message) {
	switch server.Bot().GetStatus(m.Chat.ID) {
	case fsm.BindingUsername:
		FromBindingUsernameCtrl(m)
	case fsm.BindingToken:
		FromBindingTokenCtrl(m)
	case fsm.EnablingSilent:
		FromEnablingSilentCtrl(m)
	case fsm.AddingFilter:
		FromAddingFilterCtrl(m)
	case fsm.DeletingFilter:
		FromDeletingFilterCtrl(m)
	case fsm.ActivityPage:
		FromActivityPageCtrl(m)
	case fsm.IssuePage:
		FromIssuePageCtrl(m)
	default:
		msg := fmt.Sprintf(UNKNOWN_COMMAND, m.Text)
		_ = server.Bot().Reply(m, msg)
	}
}
