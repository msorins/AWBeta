package chat

import (
	"state"
	"subscription"
)

type IChat interface {
	HandleMessages(messageReceivedCallBack func(*state.StateManager, *subscription.SubscriptionManager, string, string) []string)
	SendMessage(userId string, msgs []string)
}
