package bot

import (
	"fmt"
	"github.com/Aoi-hosizora/github-telebot/src/bot/fsm"
	"github.com/Aoi-hosizora/github-telebot/src/logger"
	"github.com/Aoi-hosizora/github-telebot/src/model"
	"github.com/Aoi-hosizora/github-telebot/src/service"
	"gopkg.in/tucnak/telebot.v2"
	"log"
	"strconv"
	"strings"
)

// noinspection GoSnakeCaseUsage
const (
	START = "This is AoiHosizora's github telebot. SendToChat /bind to bind a github user."
	HELP  = `
**Commands**
/start - show start message
/help - show this help message
/cancel - cancel the last action
/bind - bind a new github account
/unbind - unbind an old github account
/me - show the bind user information
/send - show the first page of user's events
/sendn - show the n page of user's events
/issue - show the first page of user's issue events
/issuen - show the n page of user's issue events
`
	ACTION_NO       = "There is no action now."
	ACTION_CANCELED = "Action has been canceled."
	NUM_REQUIRED    = "Excepted number, but got a string. Please resend a number."

	BIND_START         = "Please send github's username and token (split with whitespace) if you want to watch private events also. SendToChat /cancel to cancel."
	BIND_ALREADY       = "You have already bind with a github account."
	BIND_NOT_YET       = "You haven't bind a github account yet."
	BIND_FAILED        = "Failed to bind github account, please retry."
	BIND_SUCCESS       = "Bind user %s without token success. Use /send try to get events."
	BIND_TOKEN_SUCCESS = "Bind user %s with token success. Use /send try to get events."

	UNBIND_START   = "Sure to unbind the current github account %s?"
	UNBIND_FAILED  = "Failed to unbind github account, please retry."
	UNBIND_SUCCESS = "Unbind user success."

	GITHUB_ME        = "You have bind with user: %s without token."
	GITHUB_ME_TOKEN  = "You have bind with user: %s with token."
	GITHUB_FAILED    = "Failed to get github information, please retry."
	GITHUB_NOT_FOUND = "Github user not found."
	GITHUB_EMPTY     = "Empty events: \\[]"
	GITHUB_SENDN     = "Please send the page you want to get, number required."
	GITHUB_ISSUEN    = "Please send the page you want to get, number required."
)

// /start
func startCtrl(m *telebot.Message) {
	msg, err := Bot.Send(m.Chat, START)
	logger.ReplyLogger(m, msg, err)
}

// /help
func helpCtrl(m *telebot.Message) {
	msg, err := Bot.Send(m.Chat, HELP, telebot.ModeMarkdown)
	logger.ReplyLogger(m, msg, err)
}

// /cancel
func cancelCtrl(m *telebot.Message) {
	if UserStates[m.Chat.ID] == fsm.None {
		msg, err := Bot.Send(m.Chat, ACTION_NO)
		logger.ReplyLogger(m, msg, err)
	} else {
		UserStates[m.Chat.ID] = fsm.None
		msg, err := Bot.Send(m.Chat, ACTION_CANCELED)
		logger.ReplyLogger(m, msg, err)
	}
}

// btn_cancel
func cancelBtnCtrl(c *telebot.Callback) {
	// _, _ = Bot.Edit(c.Message, c.Message.Text+" (You choose Cancel)", &telebot.ReplyMarkup{})
	_ = Bot.Delete(c.Message)
	m := c.Message
	msg, err := Bot.Send(m.Chat, ACTION_CANCELED)
	logger.ReplyLogger(m, msg, err)
}

// /bind
func startBindCtrl(m *telebot.Message) {
	user := model.GetUser(m.Chat.ID)
	if user != nil {
		msg, err := Bot.Send(m.Chat, BIND_ALREADY)
		logger.ReplyLogger(m, msg, err)
	} else {
		UserStates[m.Chat.ID] = fsm.Bind
		msg, err := Bot.Send(m.Chat, BIND_START)
		logger.ReplyLogger(m, msg, err)
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
	ok, err := service.CheckUser(user.Username, user.Private, user.Token)
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
	logger.ReplyLogger(m, msg, err)
}

// /me
func meCtrl(m *telebot.Message) {
	user := model.GetUser(m.Chat.ID)
	flag := ""
	if user == nil {
		flag = BIND_NOT_YET
	} else {
		n := fmt.Sprintf("[%s](https://github.com/%s)", user.Username, user.Username)
		if user.Private {
			flag = fmt.Sprintf(GITHUB_ME_TOKEN, n)
		} else {
			flag = fmt.Sprintf(GITHUB_ME, n)
		}
	}
	msg, err := Bot.Send(m.Chat, flag)
	logger.ReplyLogger(m, msg, err)
}

// /unbind
func startUnbindCtrl(m *telebot.Message) {
	user := model.GetUser(m.Chat.ID)
	if user == nil {
		msg, err := Bot.Send(m.Chat, BIND_NOT_YET)
		logger.ReplyLogger(m, msg, err)
	} else {
		flag := fmt.Sprintf(UNBIND_START, user.Username)
		msg, err := Bot.Send(m.Chat, flag, &telebot.ReplyMarkup{
			InlineKeyboard: [][]telebot.InlineButton{{*InlineButtons["btn_unbind"]}, {*InlineButtons["btn_cancel"]}},
		})
		logger.ReplyLogger(m, msg, err)
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
	logger.ReplyLogger(m, msg, err)
}

// /send
func sendCtrl(m *telebot.Message) {
	m.Text = "1"
	sendnCtrl(m)
}

// /sendn
func startSendnCtrl(m *telebot.Message) {
	UserStates[m.Chat.ID] = fsm.Sendn
	msg, err := Bot.Send(m.Chat, GITHUB_SENDN)
	logger.ReplyLogger(m, msg, err)
}

// /sendn -> x
func sendnCtrl(m *telebot.Message) {
	page, err := strconv.Atoi(m.Text)
	if err != nil {
		msg, err := Bot.Send(m.Chat, NUM_REQUIRED)
		logger.ReplyLogger(m, msg, err)
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
		resp, err := service.GetActivityEvents(user.Username, user.Private, user.Token, page)
		if err != nil {
			log.Println(err)
			flag = GITHUB_FAILED
		} else {
			events, err := model.UnmarshalActivityEvents(resp)
			if err != nil {
				log.Println(err)
				flag = GITHUB_FAILED
			} else {
				render := service.RenderActivities(events)
				if render == "" {
					render = GITHUB_EMPTY
				} else {
					flag = fmt.Sprintf("From [%s](https://github.com/%s) (page %d):\n%s", user.Username, user.Username, page, render)
				}
			}
		}
	}

	UserStates[m.Chat.ID] = fsm.None
	msg, err := Bot.Send(m.Chat, flag, telebot.ModeMarkdown)
	logger.ReplyLogger(m, msg, err)
}

// /issue
func sendIssueCtrl(m *telebot.Message) {
	m.Text = "1"
	sendIssuenCtrl(m)
}

// /issuen
func startSendIssuenCtrl(m *telebot.Message) {
	UserStates[m.Chat.ID] = fsm.Issuen
	msg, err := Bot.Send(m.Chat, GITHUB_ISSUEN)
	logger.ReplyLogger(m, msg, err)
}

// /issuen -> x
func sendIssuenCtrl(m *telebot.Message) {
	page, err := strconv.Atoi(m.Text)
	if err != nil {
		msg, err := Bot.Send(m.Chat, NUM_REQUIRED)
		logger.ReplyLogger(m, msg, err)
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
		resp, err := service.GetIssueEvents(user.Username, user.Private, user.Token, page)
		if err != nil {
			log.Println(err)
			flag = GITHUB_FAILED
		} else {
			events, err := model.UnmarshalIssueEvents(resp)
			if err != nil {
				log.Println(err)
				flag = GITHUB_FAILED
			} else {
				render := service.RenderIssues(events)
				if render == "" {
					render = GITHUB_EMPTY
				} else {
					flag = fmt.Sprintf("From [%s](https://github.com/%s) (page %d):\n%s", user.Username, user.Username, page, render)
				}
			}
		}
	}

	UserStates[m.Chat.ID] = fsm.None
	msg, err := Bot.Send(m.Chat, flag, telebot.ModeMarkdown)
	logger.ReplyLogger(m, msg, err)
}

// onText
func onTextCtrl(m *telebot.Message) {
	switch UserStates[m.Chat.ID] {
	case fsm.Bind:
		bindCtrl(m)
	case fsm.Sendn:
		sendnCtrl(m)
	case fsm.Issuen:
		sendIssuenCtrl(m)
	default:
		msg, err := Bot.Send(m.Chat, "Unknown command: "+m.Text)
		logger.ReplyLogger(m, msg, err)
	}
}
