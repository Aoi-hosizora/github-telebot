package fsm

import (
	"github.com/Aoi-hosizora/ahlib-web/xtelebot"
)

const (
	None xtelebot.UserStatus = iota

	// Start binding, need to send username and token.
	Binding

	// Want to send activity events, need to send page number.
	ActivityPage

	// Want to send issue events, need to send page number.
	IssuePage

	// Want to setup silent, need to send 2 hour numbers.
	SilentHour
)
