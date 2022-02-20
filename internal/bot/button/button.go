package button

import (
	"github.com/Aoi-hosizora/ahlib-web/xtelebot"
)

var (
	// InlineBtnUnsubscribe (controller.Unsubscribe)
	InlineBtnUnsubscribe = xtelebot.DataBtn("Unsubscribe", "btn_unsubscribe")

	// InlineBtnCancelUnsubscribe (controller.Unsubscribe)
	InlineBtnCancelUnsubscribe = xtelebot.DataBtn("Cancel", "btn_cancel_unsubscribe")

	// InlineBtnFilter (controller.AllowIssue)
	InlineBtnFilter = xtelebot.DataBtn("Filter", "btn_filter")

	// InlineBtnNotFilter (controller.AllowIssue)
	InlineBtnNotFilter = xtelebot.DataBtn("Not Filter", "btn_not_filter")

	// InlineBtnCancelSetupIssue (controller.AllowIssue)
	InlineBtnCancelSetupIssue = xtelebot.DataBtn("Cancel", "btn_cancel_setup_issue")
)
