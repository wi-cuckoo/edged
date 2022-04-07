package protocol

// ConnectPacket MQTT Connect packet
type ConnectPacket struct {
	*Header

	ProtocolName    []byte
	ProtocolVersion byte
	CleanSession    bool
	WillFlag        bool
	WillQoS         QoSLevel
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
	// variable header
	var err error
	if c.ProtocolName, err = d.readFixedBytes(); err != nil {
		return err
	}
	if c.ProtocolVersion, err = d.readOneByte(); err != nil {
		return err
	}

	var flag byte
	if flag, err = d.readOneByte(); err != nil {
		return err
	}
	c.CleanSession = (flag>>1)&0x01 == 1
	c.WillFlag = (flag>>2)&0x01 == 1
	c.WillQoS = QoSLevel((flag >> 3) & 0x03)
	c.WillRetain = (flag>>5)&0x01 == 1
	c.PasswordFlag = (flag>>6)&0x01 == 1
	c.UsernameFlag = (flag>>7)&0x01 == 1

	if c.KeepAlive, err = d.readUint16(); err != nil {
		return err
	}

	// payload
	if c.ClientIdentifier, err = d.readFixedBytes(); err != nil {
		return err
	}
	if c.WillFlag {
		if c.WillTopic, err = d.readFixedBytes(); err != nil {
			return err
		}
		if c.WillMessage, err = d.readFixedBytes(); err != nil {
			return err
		}
	}
	if c.UsernameFlag {
		if c.Username, err = d.readFixedBytes(); err != nil {
			return err
		}
	}
	if c.PasswordFlag {
		if c.Password, err = d.readFixedBytes(); err != nil {
			return err
		}
	}

	return nil
}
