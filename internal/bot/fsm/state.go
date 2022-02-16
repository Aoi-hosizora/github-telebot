package fsm

import (
	"github.com/Aoi-hosizora/ahlib-web/xtelebot"
)

const (
	None xtelebot.ChatState = iota

	// BindingUsername (controller.Subscribe)
	BindingUsername

	// BindingToken (controller.Subscribe)
	BindingToken
)

func StateString(state xtelebot.ChatState) string {
	switch state {
	case BindingUsername, BindingToken:
		return "subscribe"
	default:
		return "?"
	}
}
