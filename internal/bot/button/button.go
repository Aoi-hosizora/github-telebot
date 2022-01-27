package button

import (
	"gopkg.in/tucnak/telebot.v2"
)

var (
	// InlineBtnUnbind (controller.Unsubscribe)
	InlineBtnUnbind = &telebot.InlineButton{Unique: "btn_unbind", Text: "Unbind"}

	// InlineBtnCancelUnbind (controller.Unsubscribe)
	InlineBtnCancelUnbind = &telebot.InlineButton{Unique: "btn_cancel_unbind", Text: "Cancel"}

	// InlineBtnFilter (controller.AllowIssue)
	InlineBtnFilter = &telebot.InlineButton{Unique: "btn_filter", Text: "Filter"}

	// InlineBtnNotFilter (controller.AllowIssue)
	InlineBtnNotFilter = &telebot.InlineButton{Unique: "btn_not_filter", Text: "Not Filter"}

	// InlineBtnCancelSetupIssue (controller.AllowIssue)
	InlineBtnCancelSetupIssue = &telebot.InlineButton{Unique: "btn_cancel_setup_issue", Text: "Cancel"}
)
