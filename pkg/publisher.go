package pkg

import "context"

// Publisher could publish a message to one topic[only]
type Publisher interface {
	Publish(ctx context.Context, m *Message) error
}
