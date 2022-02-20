package fsm

import (
	"github.com/Aoi-hosizora/ahlib-web/xtelebot"
)

const (
	None xtelebot.ChatState = iota

	// SubscribingUsername (controller.Subscribe)
	SubscribingUsername

	// SubscribingToken (controller.Subscribe)
	SubscribingToken
)

func StateString(state xtelebot.ChatState) string {
	switch state {
	case SubscribingUsername, SubscribingToken:
		return "subscribe"
	default:
		return "?"
	}
}
