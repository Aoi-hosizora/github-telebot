package fsm

type UserStatus int

const (
	None UserStatus = iota
	Binding
	Activityn
	Issuen
)
