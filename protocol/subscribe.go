package protocol

type SubscribePacket struct {
	*Header

	MessageID uint16
	Topics    [][]byte
	Qoss      []QoSLevel
}

func (s *SubscribePacket) Decode(d *decoder) error {
	var err error
	s.MessageID, err = d.readUint16()
	if err != nil {
		return err
	}
	payloadLen := s.Header.RemainLength - 2
	for payloadLen > 0 {
		topic, err := d.readFixedBytes()
		if err != nil {
			return err
		}
		s.Topics = append(s.Topics, topic)
		qos, err := d.readOneByte()
		if err != nil {
			return err
		}
		s.Qoss = append(s.Qoss, QoSLevel(qos))
		payloadLen -= 2 + len(topic) + 1 // 2 bytes of string length, plus string, plus 1 byte for Qos
	}

	return nil
}

func (s *SubscribePacket) Encode(enc *encoder) error {
	return nil
}
