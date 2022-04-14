package protocol

import (
	"bufio"
	"fmt"
	"io"
)

type Packet interface {
	MessageType() MessageType
	Decode(d *decoder) error
	Encode(e *encoder) error
}

// ReadPacket ...
func ReadPacket(r io.Reader) (Packet, error) {
	dec := &decoder{bufio.NewReader(r)}
	header, err := dec.readHeader()
	if err != nil {
		return nil, err
	}

	packet, err := selectPacket(header)
	if err != nil {
		return nil, err
	}
	err = packet.Decode(dec)
	return packet, err
}

func selectPacket(h *Header) (Packet, error) {
	switch h.MsgType {
	case CONNECT:
		return &ConnectPacket{Header: h}, nil
	case CONNACK:
		return &ConnackPacket{Header: h}, nil
	default:
		return nil, fmt.Errorf("no such MessageType: %d", h.MsgType)
	}
}
