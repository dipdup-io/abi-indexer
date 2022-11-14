package messages

import (
	"sync"

	"github.com/dipdup-net/abi-indexer/internal/random"
	"github.com/rs/zerolog/log"
)

// Subscriber -
type Subscriber struct {
	id       uint64
	topics   map[Topic]struct{}
	messages chan *Message

	mx sync.RWMutex
}

// NewSubscriber -
func NewSubscriber() (*Subscriber, error) {
	i, err := random.UInt64()
	if err != nil {
		return nil, err
	}
	return &Subscriber{
		id:       i,
		topics:   make(map[Topic]struct{}),
		messages: make(chan *Message, 1024),
	}, nil
}

// ID -
func (s *Subscriber) ID() uint64 {
	return s.id
}

// AddTopic -
func (s *Subscriber) AddTopic(topic Topic) {
	defer s.mx.Unlock()
	s.mx.Lock()

	s.topics[topic] = struct{}{}
}

// RemoveTopic -
func (s *Subscriber) RemoveTopic(topic Topic) {
	defer s.mx.Unlock()
	s.mx.Lock()

	delete(s.topics, topic)
}

// Listen -
func (s *Subscriber) Listen() <-chan *Message {
	return s.messages
}

// Notify -
func (s *Subscriber) Notify(msg *Message) {
	select {
	case s.messages <- msg:
	default:
		log.Warn().Uint64("id", s.id).Msg("can't send message: channel is full")
	}
}

// Close -
func (s *Subscriber) Close() error {
	close(s.messages)
	return nil
}
