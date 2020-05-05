package bot

import (
	"fmt"
	"github.com/Aoi-hosizora/ah-tgbot/bot/fsm"
	"github.com/Aoi-hosizora/ah-tgbot/logger"
	"github.com/Aoi-hosizora/ah-tgbot/model"
	"gopkg.in/tucnak/telebot.v2"
	"strings"
)

// /start
func startCtrl(m *telebot.Message) {
	msg, err := Bot.Send(m.Chat, "This is AoiHosizora's github telebot. Send /bind to bind a github user.")
	logger.SndLogger(m, msg, err)
}

// /cancel
func cancelCtrl(m *telebot.Message) {
	flag := ""
	if UserStates[m.Chat.ID] == fsm.None {
		flag = "There is no action now."
	} else {
		flag = "Action has been canceled."
		UserStates[m.Chat.ID] = fsm.None
	}
	msg, err := Bot.Send(m.Chat, flag)
	logger.SndLogger(m, msg, err)
}

// btn_cancel
func cancelBtnCtrl(c *telebot.Callback) {
	_, _ = Bot.Edit(c.Message, c.Message.Text+" (Cancel)")
	_, _ = Bot.EditReplyMarkup(c.Message, &telebot.ReplyMarkup{})
	m := c.Message
	msg, err := Bot.Send(m.Chat, "Action has been canceled.")
	logger.SndLogger(m, msg, err)
	c.Message.Text += " (Canceled)"
}

// /bind
func startBindCtrl(m *telebot.Message) {
	msg, err := Bot.Send(m.Chat, "Please send github's username and token (split with whitespace) if you want to watch private events also. Send /cancel to cancel.")
	logger.SndLogger(m, msg, err)
	UserStates[m.Chat.ID] = fsm.Bind
}

// /bind -> x
func bindCtrl(m *telebot.Message) {
	sp := strings.Split(m.Text, " ")
	username := sp[0]
	user := &model.User{ChatID: m.Chat.ID, Username: username}
	if len(sp) > 1 {
		user.Private = true
		user.Token = sp[1]
	}

	status := model.AddUser(user)
	flag := ""
	if status == model.DbExisted {
		flag = "You have bind with a github account, please unbind first to rebind a new account."
	} else if status == model.DbFailed {
		flag = "Failed to bind github account, database error."
	} else {
		UserStates[m.Chat.ID] = fsm.None
		flag = fmt.Sprintf("Bind user %s with success.", username)
		if user.Private == true {
			flag = fmt.Sprintf("Bind user %s with token success.", username)
		}
	}
	msg, err := Bot.Send(m.Chat, flag)
	logger.SndLogger(m, msg, err)
}

// /unbind
func startUnbindCtrl(m *telebot.Message) {
	btns := [][]telebot.InlineButton{{*InlineBtns["btn_unbind"]}, {*InlineBtns["btn_cancel"]}}
	msg, err := Bot.Send(m.Chat, "Sure to unbind the current github account?", &telebot.ReplyMarkup{
		InlineKeyboard: btns,
	})
	logger.SndLogger(m, msg, err)
}

// btn_unbind
func unbindBtnCtrl(c *telebot.Callback) {
	_, _ = Bot.Edit(c.Message, c.Message.Text+" (Unbind)")
	_, _ = Bot.EditReplyMarkup(c.Message, &telebot.ReplyMarkup{})
	m := c.Message
	status := model.DeleteUser(m.Chat.ID)
	flag := ""
	if status == model.DbNotFound {
		flag = "Account not found, maybe you don't bind a github account yet."
	} else if status == model.DbFailed {
		flag = "Failed to unbind github account, database error."
	} else {
		flag = "Unbind user success."
	}
	msg, err := Bot.Send(m.Chat, flag)
	logger.SndLogger(m, msg, err)
}

// onText
func onTextCtrl(m *telebot.Message) {
	switch UserStates[m.Chat.ID] {
	case fsm.Bind:
		bindCtrl(m)
	default:
		msg, err := Bot.Send(m.Chat, "Unknown command: "+m.Text)
		logger.SndLogger(m, msg, err)
	}
}
