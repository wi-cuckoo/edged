package pkg

import "context"

// Subscriber could subscribe message from multi topics
type Subscriber interface {
	Subscribe(ctx context.Context, m *Message) error
}
