package protocol

import "bufio"

type Decode struct {
	rb *bufio.Reader
}

func (d *Decode) readHeader() (*Header, error) {
	b, err := d.rb.ReadByte()
	if err != nil {
		return nil, err
	}

	h := new(Header)
	h.Retain = b&RetainPos == 1
	h.Qos = QoSLevel((b & QosPos) >> 1)
	h.Dup = (b&DupPos)>>3 == 1
	h.MsgType = MessageType((b & MsgTypePos) >> 4)
	return h, nil
}

func (d *Decode) readLength() (int, error) {
	bs := make([]byte, 0, 4)
	for {
		b, err := d.rb.ReadByte()
		if err != nil {
			return 0, err
		}
		bs = append(bs, b)
		// 最高位即 continuation bit 为0，表明是最后一字节
		if (b & 0xff) == 0 {
			break
		}
	}

	return calcLength(bs), nil
}

func calcLength(bs []byte) int {
	multi, value := 1, 0
	for _, b := range bs {
		value += (int(b) & 127) * multi
		multi *= 128
	}
	return value
}
