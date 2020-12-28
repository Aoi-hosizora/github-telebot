package fsm

import (
	"github.com/Aoi-hosizora/ahlib-web/xtelebot"
)

const (
	None xtelebot.UserStatus = iota

	Binding    // controller.BindCtrl -> controller.FromBindingCtrl
	SilentHour // controller.EnableSilentCtrl -> controller.FromSilentHourCtrl

	ActivityPage // controller.ActivityNCtrl -> controller.FromActivityNCtrl
	IssuePage    // controller.IssueNCtrl -> controller.FromIssueNCtrl
)
