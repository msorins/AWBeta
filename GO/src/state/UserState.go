package state

type UserState int

const (
	USER_STATE_AWB_CONFUSING            UserState = iota // == 0
	USER_STATE_AWB_OK                   UserState = iota // == 1
)
