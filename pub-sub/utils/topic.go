package utils

import (
	"sync"
)

type Topic struct {
	Name string
	Subscriber map[Subscriber]struct{}
	mu sync.RWMutex
}

func NewTopic(name string ) *Topic {
	return &Topic{
		Name: name,
		Subscriber: make(map[Subscriber]struct{}),
	}
}

func (t *Topic) AddSubscriber(sub Subscriber) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Subscriber[sub] = struct{}{}
}

func (t*Topic) RemoveSubscriber(sub Subscriber){
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.Subscriber,sub)
}

func (t*Topic) Publish(message *Message){
	t.mu.RLock()
	defer t.mu.RUnlock()
	for sub := range t.Subscriber {
		sub.OnMessage(message)
	}
}