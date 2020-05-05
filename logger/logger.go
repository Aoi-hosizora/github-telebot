package logger

import (
	"fmt"
	"github.com/Aoi-hosizora/ah-tgbot/config"
	"gopkg.in/tucnak/telebot.v2"
	"log"
)

func RcvLogger(m *telebot.Message, endpoint interface{}) {
	ep := ""
	if s, ok := endpoint.(string); ok {
		ep = s
	} else if b, ok := endpoint.(*telebot.InlineButton); ok {
		ep = b.Unique
	} else {
		ep = fmt.Sprintf("%v", endpoint)
	}
	if config.Configs.Mode == "debug" {
		if ep[0] == '\a' {
			ep = "$on_" + ep[1:]
		}
		log.Printf("[telebot] -> %4d | %18v | %d %s", m.ID, endpoint, m.Chat.ID, m.Chat.Username)
	}
}

func SndLogger(from *telebot.Message, to *telebot.Message, err error) {
	if config.Configs.Mode == "debug" {
		if from != nil {
			if err != nil {
				log.Printf("[telebot] failed to send message to %d %s: %v", from.Chat.ID, from.Chat.Username, err)
			} else if to != nil {
				du := to.Time().Sub(from.Time()).Milliseconds()
				log.Printf("[telebot] <- %4d | %6dms | -> %4d | %d %s", to.ID, du, from.ID, to.Chat.ID, to.Chat.Username)
			}
		}
	}
}
