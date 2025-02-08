package utils

import "fmt"

type PrintSub struct {
	Name string
}

func NewPrintSub (name string) *PrintSub {
	return &PrintSub{Name: name}
}

func (ps*PrintSub) OnMessage(message *Message){
	fmt.Printf("subscriber %s recived message %s\n",ps.Name,message.Content)
}