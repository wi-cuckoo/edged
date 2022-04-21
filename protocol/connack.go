package protocol

type ReturnCode byte

func (r ReturnCode) Byte() byte {
	return byte(r)
}

const (
	Accepted ReturnCode = iota
	UnacceptableProtocolVerion
	IentifierRejected
	ServerUnavailable
	AuthenticationFail
	NotAuthorized
)

type ConnackPacket struct {
	*Header

	Reserved byte
	RetCode  ReturnCode
}

func (c *ConnackPacket) Decode(d *decoder) error {
	var err error
	if c.Reserved, err = d.readOneByte(); err != nil {
		return err
	}
	var b byte
	b, err = d.readOneByte()
	c.RetCode = ReturnCode(b)
	return err
}

func (c *ConnackPacket) Encode(enc *encoder) error {
	c.Header.RemainLength = 2
	if err := enc.writeHeader(c.Header); err != nil {
		return err
	}
	_, err := enc.Write([]byte{c.Reserved, c.RetCode.Byte()})
	return err
}

func (c *ConnackPacket) MessageType() MessageType {
	return CONNACK
}
