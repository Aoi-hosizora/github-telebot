package button

import (
	"github.com/Aoi-hosizora/ahlib-web/xtelebot"
)

var (
	// InlineBtnUnbind (controller.Unsubscribe)
	InlineBtnUnbind = xtelebot.DataBtn("Unbind", "btn_unbind")

	// InlineBtnCancelUnbind (controller.Unsubscribe)
	InlineBtnCancelUnbind = xtelebot.DataBtn("Cancel", "btn_cancel_unbind")

	// InlineBtnFilter (controller.AllowIssue)
	InlineBtnFilter = xtelebot.DataBtn("Filter", "btn_filter")

	// InlineBtnNotFilter (controller.AllowIssue)
	InlineBtnNotFilter = xtelebot.DataBtn("Not Filter", "btn_not_filter")

	// InlineBtnCancelSetupIssue (controller.AllowIssue)
	InlineBtnCancelSetupIssue = xtelebot.DataBtn("Cancel", "btn_cancel_setup_issue")
)
