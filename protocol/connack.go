package protocol

type ReturnCode byte

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
