package fsm

import (
	"github.com/Aoi-hosizora/ahlib-web/xtelebot"
)

const (
	None xtelebot.ChatStatus = iota

	// controller.BindCtrl -> controller.FromBindingCtrl
	Binding

	// controller.EnableSilentCtrl -> controller.FromSilentHourCtrl
	SilentHour

	// controller.ActivityPageCtrl -> controller.FromActivityPageCtrl
	ActivityPage

	// controller.IssuePageCtrl -> controller.FromIssuePageCtrl
	IssuePage
)
