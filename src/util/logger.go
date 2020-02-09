package util

import (
	"gopkg.in/tucnak/telebot.v2"
	"log"
)

func BotLog(fromMsg *telebot.Message, sendMsg *telebot.Message, err error) {
	if fromMsg != nil { // bot
		if sendMsg == nil { // no handler
			log.Printf("[bot] receive \t |%d \t \"%s\" \t (no handler)\n", fromMsg.ID, fromMsg.Text)
		} else { // send
			if err == nil {
				timeSpan := float64(sendMsg.Time().Sub(fromMsg.Time()).Nanoseconds()) / 1e6
				log.Printf("[bot] reply \t |%d \t %.0fms \t (from %d \"%s\")\n", sendMsg.ID, timeSpan, fromMsg.ID, fromMsg.Text)
			} else {
				log.Printf("[bot] failed to reply bot of %d: %v\n", fromMsg.ID, err)
			}
		}
	} else if sendMsg != nil && sendMsg.Chat != nil { // channel
		if err == nil {
			log.Printf("[channel] send \t |%d \t \"%s\"", sendMsg.ID, sendMsg.Chat.Title)
		} else {
			log.Printf("[channel] failed to send message to channel: %v\n", err)
		}
	}
}
