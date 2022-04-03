package pkg

// Publisher could publish a message to one topic[only]
type Publisher interface {
	Publish()
}
