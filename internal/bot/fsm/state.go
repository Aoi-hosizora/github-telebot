package fsm

import (
	"github.com/Aoi-hosizora/ahlib-web/xtelebot"
)

const (
	None xtelebot.ChatStatus = iota

	// controller.BindCtrl -> controller.FromBindingUsernameCtrl
	BindingUsername

	// controller.FromBindingUsernameCtrl -> controller.FromBindingTokenCtrl
	BindingToken

	// controller.EnableSilentCtrl -> controller.FromEnablingSilentCtrl
	EnablingSilent

	// controller.AddFilterCtrl -> controller.FromAddingFilterCtrl
	AddingFilter

	// controller.DeleteFilterCtrl -> controller.FromDeletingFilterCtrl
	DeletingFilter

	// controller.ActivityPageCtrl -> controller.FromActivityPageCtrl
	ActivityPage

	// controller.IssuePageCtrl -> controller.FromIssuePageCtrl
	IssuePage
)
