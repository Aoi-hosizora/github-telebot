package controller

import (
	"github.com/Aoi-hosizora/ahlib-web/xtelebot"
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"github.com/Aoi-hosizora/github-telebot/internal/service"
	"github.com/Aoi-hosizora/github-telebot/internal/service/dao"
	"gopkg.in/tucnak/telebot.v2"
	"strings"
)

const (
	_UNEXPECTED_NUMBER = "Unexpected page number. Please send an integer value. Send /cancel to cancel."
	_EMPTY_EVENT       = "You have empty event."
)

// Activity /activity
func Activity(bw *xtelebot.BotWrapper, m *telebot.Message) {
	pageStr := strings.TrimSpace(m.Payload)
	page := 1
	if pageStr != "" {
		var err error
		page, err = xnumber.Atoi(pageStr)
		if err != nil {
			bw.ReplyTo(m, _UNEXPECTED_NUMBER)
			return
		}
		if page < 1 {
			page = 1
		}
	}

	chat, _ := dao.QueryChat(m.Chat.ID)
	if chat == nil {
		bw.ReplyTo(m, _BIND_NOT_YET)
		return
	}
	events, err := service.GetActivityEvents(chat.Username, chat.Token, page)
	if err != nil {
		bw.ReplyTo(m, _GITHUB_FAILED)
		return
	}
	formatted := service.FormatActivityEvents(events, chat.Username, page)
	if formatted == "" {
		bw.ReplyTo(m, _EMPTY_EVENT)
		return
	}
	opt := []interface{}{telebot.ModeMarkdownV2}
	if chat.Silent {
		opt = append(opt, telebot.Silent)
	}
	if !chat.Preview {
		opt = append(opt, telebot.NoPreview)
	}
	bw.ReplyTo(m, formatted, opt...)
}

// Issue /issue
func Issue(bw *xtelebot.BotWrapper, m *telebot.Message) {
	pageStr := strings.TrimSpace(m.Payload)
	page := 1
	if pageStr != "" {
		var err error
		page, err = xnumber.Atoi(pageStr)
		if err != nil {
			bw.ReplyTo(m, _UNEXPECTED_NUMBER)
			return
		}
		if page < 1 {
			page = 1
		}
	}

	chat, _ := dao.QueryChat(m.Chat.ID)
	if chat == nil {
		bw.ReplyTo(m, _BIND_NOT_YET)
		return
	}
	if chat.Token == "" {
		bw.ReplyTo(m, _ISSUE_ONLY_FOR_TOKEN)
		return
	}
	events, err := service.GetIssueEvents(chat.Username, chat.Token, page)
	if err != nil {
		bw.ReplyTo(m, _GITHUB_FAILED)
		return
	}
	if chat.FilterMe {
		events = service.FilterIssueEventSlice(events, chat.Username)
	}
	formatted := service.FormatIssueEvents(events, chat.Username, page)
	if formatted == "" {
		bw.ReplyTo(m, _EMPTY_EVENT)
		return
	}
	opt := []interface{}{telebot.ModeMarkdownV2}
	if chat.Silent {
		opt = append(opt, telebot.Silent)
	}
	if !chat.Preview {
		opt = append(opt, telebot.NoPreview)
	}
	bw.ReplyTo(m, formatted, opt...)
}
