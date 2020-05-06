package logger

import (
	"fmt"
	"github.com/Aoi-hosizora/ah-tgbot/config"
	"gopkg.in/tucnak/telebot.v2"
	"log"
)

func RcvLogger(handler interface{}, endpoint interface{}) {
	if config.Configs.Mode != "debug" {
		return
	}
	ep := ""
	if s, ok := endpoint.(string); ok {
		ep = s
	} else if b, ok := endpoint.(*telebot.InlineButton); ok {
		ep = b.Unique
	} else {
		ep = fmt.Sprintf("%v", endpoint)
	}
	if ep[0] == '\a' {
		ep = "$on_" + ep[1:]
	}

	if msg, ok := handler.(*telebot.Message); ok {
		log.Printf("[telebot] -> %4d | %18v | %d %s", msg.ID, ep, msg.Chat.ID, msg.Chat.Username)
	} else if cb, ok := handler.(*telebot.Callback); ok {
		log.Printf("[telebot] -> %4d | %18v | %d %s", cb.Message.ID, ep, cb.Message.Chat.ID, cb.Message.Chat.Username)
	} else {
		log.Printf("[telebot] -> Others | %18v", ep)
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
