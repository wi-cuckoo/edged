package internal

import (
	"context"
	"sync"
	"time"

	"github.com/wi-cuckoo/edged/protocol"
)

// topic for publish & subscribe
type topic struct {
	sync.Mutex

	name string
	done chan bool
	subs map[*subscriber]bool
}

func newtopic(name string) *topic {
	return &topic{
		name: name,
		done: make(chan bool),
		subs: make(map[*subscriber]bool),
	}
}

func (t *topic) publish(ctx context.Context, p protocol.Packet) error {
	var err error
	t.Lock()
	defer t.Unlock()
	for s := range t.subs {
		select {
		case s.ch <- p:
		case <-ctx.Done():
			return ErrCtxDone
		case <-time.After(time.Microsecond * 100):
			err = ErrTimeout
		}
	}
	return err
}

func (t *topic) subscribe(s *subscriber) {
	t.Lock()
	defer t.Unlock()
	t.subs[s] = true
}

func (t *topic) unsubscribe(s *subscriber) {
	t.Lock()
	defer t.Unlock()
	delete(t.subs, s)
}

func (t *topic) close() {
	t.Lock()
	close(t.done)
	t.Unlock()
}
