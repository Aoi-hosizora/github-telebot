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
	"strconv"
	"strings"
)

// noinspection GoSnakeCaseUsage
const (
	START           = "This is AoiHosizora's github telebot. Send /bind to bind a github user."
	ACTION_NO       = "There is no action now."
	ACTION_CANCELED = "Action has been canceled."
	NUM_REQUIRED    = "Excepted number, but got a string. Please resend a number."

	BIND_START         = "Please send github's username and token (split with whitespace) if you want to watch private events also. Send /cancel to cancel."
	BIND_ALREADY       = "You have already bind with a github account."
	BIND_NOT_YET       = "You haven't bind a github account yet."
	BIND_FAILED        = "Failed to bind github account, please retry."
	BIND_SUCCESS       = "Bind user %s without token success. Use /send try to get events."
	BIND_TOKEN_SUCCESS = "Bind user %s with token success. Use /send try to get events."

	UNBIND_START   = "Sure to unbind the current github account %s?"
	UNBIND_FAILED  = "Failed to unbind github account, please retry."
	UNBIND_SUCCESS = "Unbind user success."

	GITHUB_FAILED    = "Failed to get github information, please retry."
	GITHUB_NOT_FOUND = "Github user not found."
	GITHUB_EMPTY     = "Empty events: \\[]"
	GITHUB_SENDN     = "Please send the page you want to get, number required."
)

// /start
func startCtrl(m *telebot.Message) {
	msg, err := Bot.Send(m.Chat, START)
	logger.SndLogger(m, msg, err)
}

// /cancel
func cancelCtrl(m *telebot.Message) {
	if UserStates[m.Chat.ID] == fsm.None {
		msg, err := Bot.Send(m.Chat, ACTION_NO)
		logger.SndLogger(m, msg, err)
	} else {
		UserStates[m.Chat.ID] = fsm.None
		msg, err := Bot.Send(m.Chat, ACTION_CANCELED)
		logger.SndLogger(m, msg, err)
	}
}

// btn_cancel
func cancelBtnCtrl(c *telebot.Callback) {
	// _, _ = Bot.Edit(c.Message, c.Message.Text+" (You choose Cancel)", &telebot.ReplyMarkup{})
	_ = Bot.Delete(c.Message)
	m := c.Message
	msg, err := Bot.Send(m.Chat, ACTION_CANCELED)
	logger.SndLogger(m, msg, err)
}

// /bind
func startBindCtrl(m *telebot.Message) {
	user := model.GetUser(m.Chat.ID)
	if user != nil {
		msg, err := Bot.Send(m.Chat, BIND_ALREADY)
		logger.SndLogger(m, msg, err)
	} else {
		UserStates[m.Chat.ID] = fsm.Bind
		msg, err := Bot.Send(m.Chat, BIND_START)
		logger.SndLogger(m, msg, err)
	}
}

// /bind -> x
func bindCtrl(m *telebot.Message) {
	sp := strings.Split(m.Text, " ")
	username := sp[0]
	user := &model.User{ChatID: m.Chat.ID, Username: username}
	if len(sp) > 1 && sp[1] != "" {
		user.Private = true
		user.Token = sp[1]
	}

	flag := ""
	ok, err := util.CheckGithubUser(user.Username, user.Private, user.Token)
	if err != nil {
		log.Println(err)
		flag = GITHUB_FAILED
	} else if !ok {
		flag = GITHUB_NOT_FOUND
	} else {
		status := model.AddUser(user)
		if status == model.DbExisted {
			flag = BIND_ALREADY
		} else if status == model.DbFailed {
			flag = BIND_FAILED
		} else if user.Private {
			flag = fmt.Sprintf(BIND_TOKEN_SUCCESS, username)
		} else {
			flag = fmt.Sprintf(BIND_SUCCESS, username)
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
		msg, err := Bot.Send(m.Chat, BIND_NOT_YET)
		logger.SndLogger(m, msg, err)
	} else {
		flag := fmt.Sprintf(UNBIND_START, user.Username)
		msg, err := Bot.Send(m.Chat, flag, &telebot.ReplyMarkup{
			InlineKeyboard: [][]telebot.InlineButton{{*InlineBtns["btn_unbind"]}, {*InlineBtns["btn_cancel"]}},
		})
		logger.SndLogger(m, msg, err)
	}
}

// btn_unbind
func unbindBtnCtrl(c *telebot.Callback) {
	_ = Bot.Delete(c.Message)
	m := c.Message
	flag := ""
	status := model.DeleteUser(m.Chat.ID)
	if status == model.DbNotFound {
		flag = BIND_NOT_YET
	} else if status == model.DbFailed {
		flag = UNBIND_FAILED
	} else {
		flag = UNBIND_SUCCESS
	}
	msg, err := Bot.Send(m.Chat, flag)
	logger.SndLogger(m, msg, err)
}

// /send
func sendCtrl(m *telebot.Message) {
	m.Text = "1"
	sendnCtrl(m)
}

func startSendnCtrl(m *telebot.Message) {
	UserStates[m.Chat.ID] = fsm.Sendn
	msg, err := Bot.Send(m.Chat, GITHUB_SENDN)
	logger.SndLogger(m, msg, err)
}

// /sendn
func sendnCtrl(m *telebot.Message) {
	page, err := strconv.Atoi(m.Text)
	if err != nil {
		msg, err := Bot.Send(m.Chat, NUM_REQUIRED)
		logger.SndLogger(m, msg, err)
		return
	}
	if page <= 0 {
		page = 1
	}

	flag := ""
	user := model.GetUser(m.Chat.ID)
	if user == nil {
		flag = BIND_NOT_YET
	} else {
		reply, err := util.GetGithubEvents(user.Username, user.Private, user.Token, page)
		if err != nil {
			log.Println(err)
			flag = GITHUB_FAILED
		} else {
			events := make([]*model.GithubEvent, 0)
			err := json.Unmarshal([]byte(reply), &events)
			if err != nil {
				log.Println(err)
				flag = GITHUB_FAILED
			} else {
				flag = util.RenderGithubActions(events)
				if flag == "" {
					flag = GITHUB_EMPTY
				}
				flag = fmt.Sprintf("From [%s](https://github.com/%s) (page %d):\n%s", user.Username, user.Username, page, flag)
			}
		}
	}

	UserStates[m.Chat.ID] = fsm.None
	msg, err := Bot.Send(m.Chat, flag, telebot.ModeMarkdown)
	logger.SndLogger(m, msg, err)
}

// onText
func onTextCtrl(m *telebot.Message) {
	switch UserStates[m.Chat.ID] {
	case fsm.Bind:
		bindCtrl(m)
	case fsm.Sendn:
		sendnCtrl(m)
	default:
		msg, err := Bot.Send(m.Chat, "Unknown command: "+m.Text)
		logger.SndLogger(m, msg, err)
	}
}
