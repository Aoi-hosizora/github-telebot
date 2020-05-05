package bot

import (
	"fmt"
	"github.com/Aoi-hosizora/ah-tgbot/bot/fsm"
	"github.com/Aoi-hosizora/ah-tgbot/logger"
	"gopkg.in/tucnak/telebot.v2"
	"strings"
)

// /start
func startCtrl(m *telebot.Message) {
	msg, err := Bot.Send(m.Chat, "This is AoiHosizora's github telebot. Send /bind to bind a github user.")
	logger.SndLogger(m, msg, err)
}

// /bind
func startBindCtrl(m *telebot.Message) {
	msg, err := Bot.Send(m.Chat, "please send github's username and token (split with whitespace) if you want to watch private events also.")
	logger.SndLogger(m, msg, err)
	UserStates[m.Chat.ID] = fsm.Bind
}

// /bind -> x
func bindCtrl(m *telebot.Message) {
	sp := strings.Split(m.Text, " ")
	username := sp[0]
	flag := fmt.Sprintf("bind user %s with success.", username)
	if len(sp) > 1 {
		_ = sp[1]
		flag = fmt.Sprintf("bind user %s with token success.", username)
	}

	msg, err := Bot.Send(m.Chat, flag)
	logger.SndLogger(m, msg, err)
	UserStates[m.Chat.ID] = fsm.None
}

func onTextCtrl(m *telebot.Message) {
	switch UserStates[m.Chat.ID] {
	case fsm.Bind:
		bindCtrl(m)
	default:
		msg, err := Bot.Send(m.Chat, "Unknown command: "+m.Text)
		logger.SndLogger(m, msg, err)
	}
}
