package bot

import (
	"encoding/json"
	"fmt"
	"github.com/Aoi-hosizora/ah-tgbot/bot/fsm"
	"github.com/Aoi-hosizora/ah-tgbot/logger"
	"github.com/Aoi-hosizora/ah-tgbot/model"
	"github.com/Aoi-hosizora/ah-tgbot/util"
	"gopkg.in/tucnak/telebot.v2"
	"log"
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
	_, _ = Bot.Edit(c.Message, c.Message.Text+" (You choose Cancel)", &telebot.ReplyMarkup{})
	m := c.Message
	msg, err := Bot.Send(m.Chat, "Action has been canceled.")
	logger.SndLogger(m, msg, err)
	c.Message.Text += " (Canceled)"
}

// /bind
func startBindCtrl(m *telebot.Message) {
	user := model.GetUser(m.Chat.ID)
	if user == nil {
		UserStates[m.Chat.ID] = fsm.Bind
		msg, err := Bot.Send(m.Chat, "Please send github's username and token (split with whitespace) if you want to watch private events also. Send /cancel to cancel.")
		logger.SndLogger(m, msg, err)
	} else {
		msg, err := Bot.Send(m.Chat, "You have bind with a github account, please unbind first to rebind a new account.")
		logger.SndLogger(m, msg, err)
	}
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

	flag := ""
	ok, err := util.CheckGithubUser(user.Username, user.Private, user.Token)
	if err != nil {
		flag = "Failed to get github information, please retry."
	} else if !ok {
		flag = "User not found."
	} else {
		status := model.AddUser(user)
		if status == model.DbExisted {
			flag = "You have bind with a github account, please unbind first to rebind a new account."
		} else if status == model.DbFailed {
			flag = "Failed to bind github account, database error."
		} else {
			flag = fmt.Sprintf("Bind user %s with success.", username)
			if user.Private == true {
				flag = fmt.Sprintf("Bind user %s with token success.", username)
			}
		}
	}

	UserStates[m.Chat.ID] = fsm.None
	msg, err := Bot.Send(m.Chat, flag)
	logger.SndLogger(m, msg, err)
}

// /unbind
func startUnbindCtrl(m *telebot.Message) {
	user := model.GetUser(m.Chat.ID)
	if user == nil {
		msg, err := Bot.Send(m.Chat, "You haven't bind a github account yet.")
		logger.SndLogger(m, msg, err)
	} else {
		flag := fmt.Sprintf("Sure to unbind the current github account %s?", user.Username)
		msg, err := Bot.Send(m.Chat, flag, &telebot.ReplyMarkup{
			InlineKeyboard: [][]telebot.InlineButton{{*InlineBtns["btn_unbind"]}, {*InlineBtns["btn_cancel"]}},
		})
		logger.SndLogger(m, msg, err)
	}
}

// btn_unbind
func unbindBtnCtrl(c *telebot.Callback) {
	_, _ = Bot.Edit(c.Message, c.Message.Text+" (You choose Unbind)", &telebot.ReplyMarkup{})
	m := c.Message
	status := model.DeleteUser(m.Chat.ID)
	flag := ""
	if status == model.DbNotFound {
		flag = "You haven't bind a github account yet."
	} else if status == model.DbFailed {
		flag = "Failed to unbind github account, database error."
	} else {
		flag = "Unbind user success."
	}
	msg, err := Bot.Send(m.Chat, flag)
	logger.SndLogger(m, msg, err)
}

// /send
func sendCtrl(m *telebot.Message) {
	user := model.GetUser(m.Chat.ID)
	if user == nil {
		msg, err := Bot.Send(m.Chat, "You haven't bind with a github account, please bind first.")
		logger.SndLogger(m, msg, err)
		return
	}
	reply, err := util.GetGithubEvents(user.Username, user.Private, user.Token, 0)
	if err == nil {
		events := make([]*model.GithubEvent, 0)
		err := json.Unmarshal([]byte(reply), &events)
		if err == nil {
			render := util.RenderGithubActions(events)
			if render == "" {
				render = "Empty events: \\[]"
			}
			msg, err := Bot.Send(m.Chat, render, telebot.ModeMarkdown)
			logger.SndLogger(m, msg, err)
			return
		} else {
			log.Println(err)
		}
	} else {
		log.Println(err)
	}

	msg, err := Bot.Send(m.Chat, "Get github reply failed, please retry.")
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
