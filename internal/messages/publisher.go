package messages

import "sync"

// Subscribers -
type Subscribers map[uint64]*Subscriber

// Publisher -
type Publisher struct {
	subscibers Subscribers
	topics     map[Topic]Subscribers

	mx sync.RWMutex
}

// NewPublisher -
func NewPublisher() *Publisher {
	return &Publisher{
		subscibers: make(Subscribers),
		topics:     make(map[Topic]Subscribers),
	}
}

// Notify -
func (publisher *Publisher) Notify(msg *Message) {
	if msg == nil {
		return
	}

	defer publisher.mx.RUnlock()
	publisher.mx.RLock()

	if subscribers, ok := publisher.topics[msg.topic]; ok {
		for _, subscriber := range subscribers {
			// TODO: non-blocking message sendning with case when subscriber is stucked
			subscriber.Notify(msg)
		}
	}
}

// Subscribe -
func (publisher *Publisher) Subscribe(subscriber *Subscriber, topic Topic) {
	if subscriber == nil {
		return
	}

	defer publisher.mx.Unlock()
	publisher.mx.Lock()

	if _, ok := publisher.topics[topic]; !ok {
		publisher.topics[topic] = make(Subscribers)
	}
	subscriber.AddTopic(topic)
	publisher.topics[topic][subscriber.id] = subscriber
}

// Unsubscribe -
func (publisher *Publisher) Unsubscribe(subscriber *Subscriber, topic Topic) {
	if subscriber == nil {
		return
	}

	defer publisher.mx.Unlock()
	publisher.mx.Lock()

	if subscribers, ok := publisher.topics[topic]; ok {
		delete(subscribers, subscriber.id)
	}

	subscriber.RemoveTopic(topic)
}
