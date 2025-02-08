package main

import (
	"pub-sub/utils"
)

func Run() {
	topic1 := utils.NewTopic("Topic1")
	topic2 :=utils.NewTopic("Topic2")

	// Create publishers
	publisher1 :=utils.NewPublisher()
	publisher2 := utils.NewPublisher()

	// Create subscribers
	subscriber1 := utils.NewPrintSub("Subscriber1")
	subscriber2 := utils.NewPrintSub("Subscriber2")
	subscriber3 := utils.NewPrintSub("Subscriber3")

	publisher1.RegisterTopic(topic1)
	publisher2.RegisterTopic(topic2)

	// Subscribe to topics
	topic1.AddSubscriber(subscriber1)
	topic1.AddSubscriber(subscriber2)
	topic2.AddSubscriber(subscriber2)
	topic2.AddSubscriber(subscriber3)

	// Publish messages
	publisher1.Publish(topic1, utils.NewMessage("Message1 for Topic1"))
	publisher1.Publish(topic1, utils.NewMessage("Message2 for Topic1"))
	publisher2.Publish(topic2, utils.NewMessage("Message1 for Topic2"))

	// Unsubscribe from a topic
	topic1.RemoveSubscriber(subscriber2)

	// Publish more messages
	publisher1.Publish(topic1, utils.NewMessage("Message3 for Topic1"))
	publisher2.Publish(topic2, utils.NewMessage("Message2 for Topic2"))
}

func main(){
	Run()
}