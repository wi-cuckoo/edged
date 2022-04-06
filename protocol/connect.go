package protocol

// ConnectPacket MQTT Connect packet
type ConnectPacket struct {
	Header

	ProtocolName    []byte
	ProtocolVersion byte
	CleanSession    bool
	WillFlag        bool
	WillQos         byte
	WillRetain      bool
	UsernameFlag    bool
	PasswordFlag    bool
	ReservedBit     byte
	KeepAlive       uint16 // second

	ClientIdentifier []byte
	WillTopic        []byte
	WillMessage      []byte
	Username         []byte
	Password         []byte
}

func (c *ConnectPacket) Decode(d *decoder) error {
	// protocol name
	var err error
	if c.ProtocolName, err = d.readFixedBytes(); err != nil {
		return err
	}
	if c.ProtocolVersion, err = d.ReadByte(); err != nil {
		return err
	}

	var flag byte
	if flag, err = d.ReadByte(); err != nil {
		return err
	}
	c.CleanSession = (flag>>1)&0x01 == 1
	// todo

	if c.KeepAlive, err = d.readUint16(); err != nil {
		return err
	}
	if c.ClientIdentifier, err = d.readFixedBytes(); err != nil {
		return err
	}

	return err
}
