package internal

import (
	"context"
	"errors"
	"sync"

	"github.com/wi-cuckoo/edged/protocol"
)

type pubsub struct {
	sync.RWMutex

	topics map[string]*topic
}

func (ps *pubsub) createTopic(name string) *topic {
	ps.Lock()
	t, ok := ps.topics[name]
	if !ok {
		ps.topics[name] = newtopic(name)
	}
	ps.Unlock()
	return t
}

func (ps *pubsub) publish(ctx context.Context, topicName string, p protocol.Packet) error {
	ps.RLock()
	t, ok := ps.topics[topicName]
	ps.RUnlock()
	if !ok {
		return ErrNotFound
	}
	return t.publish(ctx, p)
}

func (ps *pubsub) subscribe(ctx context.Context, topicName string) (<-chan protocol.Packet, error) {
	ps.RLock()
	t, ok := ps.topics[topicName]
	ps.RUnlock()
	if !ok {
		return nil, ErrNotFound
	}

	s := &subscriber{
		ch: make(chan protocol.Packet, 1<<7),
	}
	t.subscribe(s)

	return s.ch, nil
}

type subscriber struct {
	ch chan protocol.Packet
}

// ErrNotFound is returned when the named topic does not exist.
var ErrNotFound = errors.New("topic not found")

// ErrTimeout is returned when publish timeout
var ErrTimeout = errors.New("publish timeout")

// ErrCtxDone is returned when publish timeout
var ErrCtxDone = errors.New("context done")

// ErrExisted ....
var ErrExisted = errors.New("portal has existed")

// ErrFullChan ....
var ErrFullChan = errors.New("portal channel is full")

// ErrNoReceiver ....
var ErrNoReceiver = errors.New("portal receiver has gone")
