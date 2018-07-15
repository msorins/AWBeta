package chat

type IChat interface {
	HandleMessages(messageReceivedCallBack func(string, string) []string)
	SendMessage(userId string, msgs []string)
}
