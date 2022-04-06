package internal

import (
	"context"

	"github.com/wi-cuckoo/edged/pkg"
)

// Topic for publish & subscribe
type Topic struct {
	Name string

	retainMsg *pkg.Message
	subs      []pkg.Subscriber
}

func (t *Topic) RegisterSub(sub pkg.Subscriber) {
	t.subs = append(t.subs, sub)
	if t.retainMsg != nil {
		sub.Subscribe(context.TODO(), t.retainMsg)
	}
}

func (t *Topic) Process(msg *pkg.Message) {

}
