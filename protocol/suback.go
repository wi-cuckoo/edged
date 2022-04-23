package protocol

type SubackPacket struct {
	*Header
	MessageID   uint16
	ReturnCodes []ReturnCode
}

func (s *SubackPacket) Encode(enc *encoder) error {
	s.Header.RemainLength = 2 + len(s.ReturnCodes)
	if err := enc.writeHeader(s.Header); err != nil {
		return err
	}
	enc.writeUint16(s.MessageID)
	for i := 0; i < len(s.ReturnCodes); i++ {
		enc.writeOneByte(byte(s.ReturnCodes[i]))
	}
	return nil
}

func (s *SubackPacket) Decode(d *decoder) error {
	return nil
}
