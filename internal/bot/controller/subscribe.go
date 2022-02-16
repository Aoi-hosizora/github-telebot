package controller

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/xtelebot"
	"github.com/Aoi-hosizora/ahlib/xstatus"
	"github.com/Aoi-hosizora/ahlib/xstring"
	"github.com/Aoi-hosizora/github-telebot/internal/bot/button"
	"github.com/Aoi-hosizora/github-telebot/internal/bot/fsm"
	"github.com/Aoi-hosizora/github-telebot/internal/pkg/config"
	"github.com/Aoi-hosizora/github-telebot/internal/service"
	"github.com/Aoi-hosizora/github-telebot/internal/service/dao"
	"gopkg.in/tucnak/telebot.v2"
	"strings"
)

const (
	_BIND_ALREADY        = "You have already bind with a github account."
	_BIND_USERNAME_Q     = "Please send the github's username which you want to bind. Send /cancel to cancel."
	_BIND_TOKEN_Q        = "Do you want to watch private events? Send your token if you want, otherwise send 'no'. Send /cancel to cancel."
	_BIND_EMPTY_USERNAME = "Please send a non-empty username (without whitespace). Send /cancel to cancel."
	_BIND_EMPTY_TOKEN    = "Please send a non-empty token (without whitespace), or send 'no' to ignore. Send /cancel to cancel."
	_BIND_FAILED         = "Failed to bind github account, please retry later."
	_BIND_TOKEN_SUCCESS  = "Binding user '%s' with token success. " + _BIND_SUCCESS_TIP
	_BIND_NOTOK_SUCCESS  = "Binding user '%s' without token success. " + _BIND_SUCCESS_TIP
	_BIND_SUCCESS_TIP    = "Send /activity to get activity events, send /issue to get issue events.\n" +
		"(Tips: new activity events will be sent periodically, but issue events will not be sent in default, you can use /allowissue to allow sending periodically)"

	_BIND_NOT_YET   = "You have not bind with a github account yet."
	_UNBIND_Q       = "Sure to unbind the current github account '%s'?"
	_UNBIND_FAILED  = "Failed to unbind github account, please retry later."
	_UNBIND_SUCCESS = "Unbind user success."

	_GITHUB_FAILED         = "Failed to query information from github, please retry later."
	_GITHUB_USER_NOT_FOUND = "Github user is not found, or the token is invalid."
	_GITHUB_ME_NOTOK       = "You have bound with user '%s' without token, current options: %s"
	_GITHUB_ME_TOKEN       = "You have bound with user '%s' with token '%s', current options: %s"
)

// Subscribe /subscribe
func Subscribe(bw *xtelebot.BotWrapper, m *telebot.Message) {
	{
		chat, _ := dao.QueryChat(m.Chat.ID)
		if chat != nil {
			bw.RespondReply(m, false, _BIND_ALREADY)
		} else {
			bw.Data().SetState(m.Chat.ID, fsm.BindingUsername)
			bw.RespondReply(m, false, _BIND_USERNAME_Q)
		}
	}

	const usernameKey = "username"
	if !bw.Shs().IsRegistered(fsm.BindingUsername) {
		bw.Shs().Register(fsm.BindingUsername, func(bw *xtelebot.BotWrapper, m *telebot.Message) {
			username := strings.TrimSpace(m.Text)
			if username == "" {
				bw.RespondReply(m, false, _BIND_EMPTY_USERNAME)
			} else {
				bw.Data().SetCache(m.Chat.ID, usernameKey, username)
				bw.Data().SetState(m.Chat.ID, fsm.BindingToken)
				bw.RespondReply(m, false, _BIND_TOKEN_Q)
			}
		})
	}
	if !bw.Shs().IsRegistered(fsm.BindingToken) {
		bw.Shs().Register(fsm.BindingToken, func(bw *xtelebot.BotWrapper, m *telebot.Message) {
			token := strings.TrimSpace(m.Text)
			if token == "" {
				bw.RespondReply(m, false, _BIND_EMPTY_TOKEN)
				return
			}
			v, _ := bw.Data().GetCache(m.Chat.ID, usernameKey)
			username := v.(string)
			bw.Data().RemoveCache(m.Chat.ID, usernameKey)
			if strings.ToLower(token) == "no" {
				token = ""
			}

			flag := ""
			ok, err := service.CheckUserExistence(username, token)
			if err != nil {
				flag = _GITHUB_FAILED
			} else if !ok {
				flag = _GITHUB_USER_NOT_FOUND
			} else {
				status, err := dao.CreateChat(m.Chat.ID, username, token)
				if status == xstatus.DbExisted {
					flag = _BIND_ALREADY
				} else if status == xstatus.DbFailed {
					flag = _BIND_FAILED
					if config.IsDebugMode() {
						flag += "\nError：" + err.Error()
					}
				} else if token != "" {
					flag = fmt.Sprintf(_BIND_TOKEN_SUCCESS, username)
				} else {
					flag = fmt.Sprintf(_BIND_NOTOK_SUCCESS, username)
				}
			}
			bw.Data().DeleteState(m.Chat.ID)
			bw.RespondReply(m, false, flag)
		})
	}
}

// Unsubscribe /unsubscribe
func Unsubscribe(bw *xtelebot.BotWrapper, m *telebot.Message) {
	{
		chat, _ := dao.QueryChat(m.Chat.ID)
		if chat == nil {
			bw.RespondReply(m, false, _BIND_NOT_YET)
		} else {
			bw.RespondReply(m, false, fmt.Sprintf(_UNBIND_Q, chat.Username), xtelebot.SetInlineKeyboard(xtelebot.InlineKeyboard(
				xtelebot.InlineRow{button.InlineBtnUnbind},
				xtelebot.InlineRow{button.InlineBtnCancelUnbind},
			)))
		}
	}

	if !bw.IsHandled(button.InlineBtnUnbind) {
		bw.HandleInlineButton(button.InlineBtnUnbind, func(bw *xtelebot.BotWrapper, c *telebot.Callback) {
			flag := ""
			status, err := dao.DeleteChat(m.Chat.ID)
			if status == xstatus.DbNotFound {
				flag = _BIND_NOT_YET
			} else if status == xstatus.DbFailed {
				flag = _UNBIND_FAILED
				if config.IsDebugMode() {
					flag += "\nError：" + err.Error()
				}
			} else {
				flag = _UNBIND_SUCCESS
			}
			bw.Bot().Edit(c.Message, flag, xtelebot.RemoveInlineKeyboard())
		})
	}
	if !bw.IsHandled(button.InlineBtnCancelUnbind) {
		bw.HandleInlineButton(button.InlineBtnCancelUnbind, func(bw *xtelebot.BotWrapper, c *telebot.Callback) {
			bw.Bot().Edit(c.Message, fmt.Sprintf("%s (canceled)", c.Message.Text), xtelebot.RemoveInlineKeyboard())
		})
	}
}

// Me /me
func Me(bw *xtelebot.BotWrapper, m *telebot.Message) {
	flag := ""
	chat, _ := dao.QueryChat(m.Chat.ID)
	if chat == nil {
		flag = _BIND_NOT_YET
	} else {
		url := fmt.Sprintf("[%s](https://github.com/%s)", chat.Username, chat.Username)
		if chat.Token == "" {
			flag = fmt.Sprintf(_GITHUB_ME_NOTOK, url)
		} else {
			tok := xstring.MaskTokenR(chat.Token, '*', 0, 1, 2, -1, -2, -3)
			tok = strings.ReplaceAll(tok, "*", "\\*")
			flag = fmt.Sprintf(_GITHUB_ME_TOKEN, url, tok)
		}
	}
	bw.RespondReply(m, false, flag, telebot.ModeMarkdown)
}
