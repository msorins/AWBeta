package wit

type MessageIntent int

const (
	MESSAGE_NO_INTENT              MessageIntent  = iota // == 0
	MESSAGE_REQUEST_ALL_HISTORY     MessageIntent = iota // == 1
	MESSAGE_REQUEST_SUBSCRIPTION	MessageIntent = iota // == 2
)
