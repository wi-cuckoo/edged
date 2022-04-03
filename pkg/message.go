package pkg

type Message struct {
	TopName string
	data    []byte
	Retain  bool
}

func (m *Message) CopyData() []byte {
	cp := make([]byte, len(m.data))
	copy(cp, m.data)
	return cp
}
