package bot

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/xstatus"
	"github.com/Aoi-hosizora/github-telebot/src/bot/fsm"
	"github.com/Aoi-hosizora/github-telebot/src/model"
	"github.com/Aoi-hosizora/github-telebot/src/service"
	"gopkg.in/tucnak/telebot.v2"
	"strings"
)

// /bind
func BindCtrl(m *telebot.Message) {
	user := model.GetUser(m.Chat.ID)
	if user != nil {
		_ = Bot.Reply(m, BIND_ALREADY)
	} else {
		Bot.UserStates[m.Chat.ID] = fsm.Binding
		_ = Bot.Reply(m, BIND_START)
	}
}

// /bind -> x
func fromBindingCtrl(m *telebot.Message) {
	text := strings.TrimSpace(m.Text)
	if text == "" {
		_ = Bot.Reply(m, BIND_EMPTY)
		return
	}

	sp := strings.Split(text, " ")
	username := sp[0]
	user := &model.User{ChatID: m.Chat.ID, Username: username}
	if len(sp) >= 2 {
		user.Token = sp[1]
	}

	flag := ""
	ok, err := service.CheckUser(user.Username, user.Token)
	if err != nil {
		flag = GITHUB_FAILED
	} else if !ok {
		flag = GITHUB_NOT_FOUND
	} else {
		status := model.AddUser(user)
		if status == xstatus.DbExisted {
			flag = BIND_ALREADY
		} else if status == xstatus.DbFailed {
			flag = BIND_FAILED
		} else if user.Token != "" {
			flag = fmt.Sprintf(BIND_TOKEN_SUCCESS, username)
		} else {
			flag = fmt.Sprintf(BIND_SUCCESS, username)
		}
	}

	Bot.UserStates[m.Chat.ID] = fsm.None
	_ = Bot.Reply(m, flag)
}

// /me
func MeCtrl(m *telebot.Message) {
	user := model.GetUser(m.Chat.ID)
	flag := ""
	if user == nil {
		flag = BIND_NOT_YET
	} else {
		n := fmt.Sprintf("[%s](https://github.com/%s)", user.Username, user.Username)
		if user.Token != "" {
			flag = fmt.Sprintf(GITHUB_ME_TOKEN, n)
		} else {
			flag = fmt.Sprintf(GITHUB_ME, n)
		}
	}
	_ = Bot.Reply(m, flag)
}

// /unbind
func UnbindCtrl(m *telebot.Message) {
	user := model.GetUser(m.Chat.ID)
	if user == nil {
		_ = Bot.Reply(m, BIND_NOT_YET)
		return
	}

	flag := fmt.Sprintf(UNBIND_START, user.Username)
	_ = Bot.Reply(m, flag, &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{
			{*Bot.InlineButtons["btn_unbind"]}, {*Bot.InlineButtons["btn_cancel"]},
		},
	})
}

// inl:btn_cancel
func InlBtnCancelCtrl(c *telebot.Callback) {
	_ = Bot.Delete(c.Message)
	_ = Bot.Reply(c.Message, ACTION_CANCELED)
}

// inl:btn_unbind
func InlBtnUnbindCtrl(c *telebot.Callback) {
	_ = Bot.Delete(c.Message)
	flag := ""
	status := model.DeleteUser(c.Message.Chat.ID)
	if status == xstatus.DbNotFound {
		flag = BIND_NOT_YET
	} else if status == xstatus.DbFailed {
		flag = UNBIND_FAILED
	} else {
		flag = UNBIND_SUCCESS
	}

	_ = Bot.Reply(c.Message, flag)
}
