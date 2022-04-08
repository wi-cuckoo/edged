package protocol

import (
	"errors"
	"io"
)

type PublishPacket struct {
	*Header

	TopicName []byte
	MessageID uint16
	Payload   []byte
}

func (p *PublishPacket) Decode(d *decoder) error {
	var err error
	payloadLength := p.RemainLength
	if p.TopicName, err = d.readFixedBytes(); err != nil {
		return err
	}

	payloadLength -= len(p.TopicName) + 2
	if p.Qos > QoSLevel0 {
		if p.MessageID, err = d.readUint16(); err != nil {
			return err
		}
		payloadLength -= 2
	}

	if payloadLength < 0 {
		return errors.New("negative payload length")
	}

	p.Payload = make([]byte, payloadLength)
	_, err = io.ReadFull(d, p.Payload)
	return err
}
