package controller

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib-web/xtelebot"
	"github.com/Aoi-hosizora/ahlib/xgeneric/xsugar"
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
	_SUBSCRIBE_ALREADY        = "You have already subscribed with a GitHub account."
	_SUBSCRIBE_USERNAME_Q     = "Please send your GitHub username. Send /cancel to cancel the current action."
	_SUBSCRIBE_TOKEN_Q        = "Do you want to watch private events? Send your token with repo scope if you want, otherwise send 'no'. Send /cancel to cancel the current action."
	_SUBSCRIBE_EMPTY_USERNAME = "Please send a non-empty username (without whitespace). Send /cancel to cancel the current action."
	_SUBSCRIBE_EMPTY_TOKEN    = "Please send a non-empty token (without whitespace), or send 'no' for no token. Send /cancel to cancel the current action."
	_SUBSCRIBE_FAILED         = "Oops. Failed to subscribe GitHub account, please retry later."
	_SUBSCRIBE_TOKEN_SUCCESS  = "Done. Subscribe GitHub account '%s' with token successfully. " + _SUBSCRIBE_SUCCESS_TIP
	_SUBSCRIBE_NOTOK_SUCCESS  = "Done. Subscribe GitHub account '%s' without token successfully. " + _SUBSCRIBE_SUCCESS_TIP
	_SUBSCRIBE_SUCCESS_TIP    = "Now you can send /activity to get activity events, and send /issue to get issue events if you subscribe with token.\n" +
		"(Tips: new activity events will be notified periodically, but issue events will not in default, you can send /allowissue to enable this action)"

	_SUBSCRIBE_NOT_YET   = "You have not subscribed with any GitHub account yet."
	_UNSUBSCRIBE_Q       = "Sure to unsubscribe the current GitHub account '%s'?"
	_UNSUBSCRIBE_FAILED  = "Oops. Failed to unsubscribe GitHub account, please retry later."
	_UNSUBSCRIBE_SUCCESS = "Done. Unsubscribe successfully."

	_GITHUB_FAILED         = "Oops. Failed to fetch information from GitHub, please retry later."
	_GITHUB_USER_NOT_FOUND = "Oops. Github user you specified is not found, or the token you sent is invalid."
	_GITHUB_ME             = "You have subscribed with account '%s', current options:\n%s"
)

// Subscribe /subscribe
func Subscribe(bw *xtelebot.BotWrapper, m *telebot.Message) {
	{
		chat, _ := dao.QueryChat(m.Chat.ID)
		if chat != nil {
			bw.RespondReply(m, false, _SUBSCRIBE_ALREADY)
		} else {
			bw.Data().SetState(m.Chat.ID, fsm.SubscribingUsername)
			bw.RespondReply(m, false, _SUBSCRIBE_USERNAME_Q)
		}
	}

	const usernameKey = "username"
	if !bw.Shs().IsRegistered(fsm.SubscribingUsername) {
		bw.Shs().Register(fsm.SubscribingUsername, func(bw *xtelebot.BotWrapper, m *telebot.Message) {
			username := strings.TrimSpace(m.Text)
			if username == "" {
				bw.RespondReply(m, false, _SUBSCRIBE_EMPTY_USERNAME)
			} else {
				bw.Data().SetCache(m.Chat.ID, usernameKey, username)
				bw.Data().SetState(m.Chat.ID, fsm.SubscribingToken)
				bw.RespondReply(m, false, _SUBSCRIBE_TOKEN_Q)
			}
		})
	}
	if !bw.Shs().IsRegistered(fsm.SubscribingToken) {
		bw.Shs().Register(fsm.SubscribingToken, func(bw *xtelebot.BotWrapper, m *telebot.Message) {
			token := strings.TrimSpace(m.Text)
			if token == "" {
				bw.RespondReply(m, false, _SUBSCRIBE_EMPTY_TOKEN)
				return
			}
			v, _ := bw.Data().GetCache(m.Chat.ID, usernameKey)
			bw.Data().RemoveCache(m.Chat.ID, usernameKey)
			username := v.(string)
			if strings.ToLower(token) == "no" {
				token = ""
			}

			// <<<<<<
			flag := ""
			ok, err := service.CheckUserExistence(username, token)
			if err != nil {
				flag = _GITHUB_FAILED
			} else if !ok {
				flag = _GITHUB_USER_NOT_FOUND
			} else {
				status, err := dao.CreateChat(m.Chat.ID, username, token)
				if status == xstatus.DbExisted {
					flag = _SUBSCRIBE_ALREADY
				} else if status == xstatus.DbFailed {
					flag = _SUBSCRIBE_FAILED
					if config.IsDebugMode() {
						flag += "\nError：" + err.Error()
					}
				} else if token != "" {
					flag = fmt.Sprintf(_SUBSCRIBE_TOKEN_SUCCESS, username)
				} else {
					flag = fmt.Sprintf(_SUBSCRIBE_NOTOK_SUCCESS, username)
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
			bw.RespondReply(m, false, _SUBSCRIBE_NOT_YET)
		} else {
			bw.RespondReply(m, false, fmt.Sprintf(_UNSUBSCRIBE_Q, chat.Username), xtelebot.SetInlineKeyboard(xtelebot.InlineKeyboard(
				xtelebot.InlineRow{button.InlineBtnUnsubscribe},
				xtelebot.InlineRow{button.InlineBtnCancelUnsubscribe},
			)))
		}
	}

	if !bw.IsHandled(button.InlineBtnUnsubscribe) {
		bw.HandleInlineButton(button.InlineBtnUnsubscribe, func(bw *xtelebot.BotWrapper, c *telebot.Callback) {
			flag := ""
			status, err := dao.DeleteChat(m.Chat.ID)
			if status == xstatus.DbNotFound {
				flag = _SUBSCRIBE_NOT_YET
			} else if status == xstatus.DbFailed {
				flag = _UNSUBSCRIBE_FAILED
				if config.IsDebugMode() {
					flag += "\nError：" + err.Error()
				}
			} else {
				flag = _UNSUBSCRIBE_SUCCESS
			}
			bw.RespondEdit(c.Message, flag, xtelebot.RemoveInlineKeyboard())
			bw.RespondCallback(c, nil)
		})
	}
	if !bw.IsHandled(button.InlineBtnCancelUnsubscribe) {
		bw.HandleInlineButton(button.InlineBtnCancelUnsubscribe, func(bw *xtelebot.BotWrapper, c *telebot.Callback) {
			bw.RespondEdit(c.Message, fmt.Sprintf("%s (canceled)", c.Message.Text), xtelebot.RemoveInlineKeyboard())
			bw.RespondCallback(c, nil)
		})
	}
}

// Me /me
func Me(bw *xtelebot.BotWrapper, m *telebot.Message) {
	flag := ""
	chat, _ := dao.QueryChat(m.Chat.ID)
	if chat == nil {
		flag = _SUBSCRIBE_NOT_YET
	} else {
		options := make([]string, 0, 4)
		if chat.Token == "" {
			options = append(options, "subscribe without token")
		} else {
			masked := xstring.StringMaskTokenR(chat.Token, "\\*", 0, 1, 2, 3, -1, -2, -3, -4)
			options = append(options, fmt.Sprintf("subscribe with token \"%s\"", masked))
		}
		options = append(options, xsugar.IfThenElse(!chat.Issue, "disallow to notify new issue events",
			xsugar.IfThenElse(chat.FilterMe, "notify new issue events excluding from me", "notify new issue events including from me")))
		options = append(options, xsugar.IfThenElse(chat.Silent, "send message with no notification", "send message with notification"))
		options = append(options, xsugar.IfThenElse(chat.Preview, "enable preview for links", "disable preview for links"))
		bs := strings.Builder{}
		for i, opt := range options {
			if bs.Len() > 0 {
				bs.WriteByte('\n')
			}
			bs.WriteString(fmt.Sprintf("%d. %s", i+1, opt))
		}
		url := fmt.Sprintf("[%s](https://github.com/%s)", chat.Username, chat.Username)
		flag = fmt.Sprintf(_GITHUB_ME, url, bs.String())
	}
	bw.RespondReply(m, false, flag, telebot.ModeMarkdown)
}
