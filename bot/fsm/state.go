package fsm

type UserStatus int

const (
	None UserStatus = iota
	Bind
	Sendn
	Issuen
)
