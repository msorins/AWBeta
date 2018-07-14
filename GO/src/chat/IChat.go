package chat

import "state"

type IChat interface {
	HandleMessages(messageReceivedCallBack func(*state.StateManager, string, string) []string)
	SendMessage(userId string, msgs []string)
}
