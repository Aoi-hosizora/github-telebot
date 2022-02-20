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
	_INVALID_NUMBER = "Expected page number, but got an invalid parameter. Please send an integer value. Send /cancel to cancel the current action."
	_EMPTY_EVENT    = "Currently you have empty event."
)

// Activity /activity
func Activity(bw *xtelebot.BotWrapper, m *telebot.Message) {
	pageStr := strings.TrimSpace(m.Payload)
	page := 1
	if pageStr != "" {
		var err error
		page, err = xnumber.Atoi(pageStr)
		if err != nil {
			bw.RespondReply(m, false, _INVALID_NUMBER)
			return
		}
		if page < 1 {
			page = 1
		}
	}

	chat, _ := dao.QueryChat(m.Chat.ID)
	if chat == nil {
		bw.RespondReply(m, false, _SUBSCRIBE_NOT_YET)
		return
	}
	events, err := service.GetActivityEvents(chat.Username, chat.Token, page)
	if err != nil {
		bw.RespondReply(m, false, _GITHUB_FAILED)
		return
	}
	formatted := service.FormatActivityEvents(events, chat.Username, page) // <<< MarkdownV2
	if formatted == "" {
		bw.RespondReply(m, false, _EMPTY_EVENT)
		return
	}
	opts := []interface{}{telebot.ModeMarkdownV2}
	if chat.Silent {
		opts = append(opts, telebot.Silent)
	}
	if !chat.Preview {
		opts = append(opts, telebot.NoPreview)
	}
	bw.RespondReply(m, false, formatted, opts...)
}

// Issue /issue
func Issue(bw *xtelebot.BotWrapper, m *telebot.Message) {
	pageStr := strings.TrimSpace(m.Payload)
	page := 1
	if pageStr != "" {
		var err error
		page, err = xnumber.Atoi(pageStr)
		if err != nil {
			bw.RespondReply(m, false, _INVALID_NUMBER)
			return
		}
		if page < 1 {
			page = 1
		}
	}

	chat, _ := dao.QueryChat(m.Chat.ID)
	if chat == nil {
		bw.RespondReply(m, false, _SUBSCRIBE_NOT_YET)
		return
	}
	if chat.Token == "" {
		bw.RespondReply(m, false, _ISSUE_ONLY_FOR_TOKEN)
		return
	}
	events, err := service.GetIssueEvents(chat.Username, chat.Token, page)
	if err != nil {
		bw.RespondReply(m, false, _GITHUB_FAILED)
		return
	}
	if chat.FilterMe {
		events = service.FilterIssueEventSlice(events, chat.Username)
	}
	formatted := service.FormatIssueEvents(events, chat.Username, page) // <<< MarkdownV2
	if formatted == "" {
		bw.RespondReply(m, false, _EMPTY_EVENT)
		return
	}
	opts := []interface{}{telebot.ModeMarkdownV2}
	if chat.Silent {
		opts = append(opts, telebot.Silent)
	}
	if !chat.Preview {
		opts = append(opts, telebot.NoPreview)
	}
	bw.RespondReply(m, false, formatted, opts...)
}
