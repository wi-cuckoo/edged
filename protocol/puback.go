package protocol

type PubackPacket struct {
	*Header

	MessageID uint16
}

func (p *PubackPacket) Decode(d *decoder) error {
	var err error
	p.MessageID, err = d.readUint16()
	return err
}

func (p *PubackPacket) Encode(enc *encoder) error {
	var err error
	p.Header.RemainLength = 2
	if err = enc.writeHeader(p.Header); err != nil {
		return err
	}
	return enc.writeUint16(p.MessageID)
}
