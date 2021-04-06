package button

import (
	"gopkg.in/tucnak/telebot.v2"
)

var (
	// controller.UnbindCtrl & controller.AllowIssueCtrl -> controller.InlBtnCancelCtrl
	InlineBtnCancel = &telebot.InlineButton{Unique: "btn_cancel", Text: "Cancel"}

	// controller.UnbindCtrl -> controller.InlBtnUnbindCtrl
	InlineBtnUnbind = &telebot.InlineButton{Unique: "btn_unbind", Text: "Unbind"}

	// controller.AllowIssueCtrl -> controller.InlBtnFilterCtrl
	InlineBtnFilter = &telebot.InlineButton{Unique: "btn_filter", Text: "Filter"}

	// controller.AllowIssueCtrl -> controller.InlBtnNotFilterCtrl
	InlineBtnNotFilter = &telebot.InlineButton{Unique: "btn_not_filter", Text: "Not Filter"}
)
