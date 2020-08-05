package fsm

import (
	"github.com/Aoi-hosizora/ahlib-web/xtelebot"
)

const (
	None xtelebot.UserStatus = iota
	Binding
	ActivityN
	IssueN
)
