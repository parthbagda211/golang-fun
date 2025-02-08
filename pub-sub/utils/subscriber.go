package utils

type Subscriber interface {
	OnMessage(message *Message)
}