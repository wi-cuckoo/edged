package protocol

import (
	"bufio"
	"encoding/binary"
	"io"
)

type decoder struct {
	*bufio.Reader
}

func (d *decoder) readHeader() (*Header, error) {
	b, err := d.ReadByte()
	if err != nil {
		return nil, err
	}

	h := new(Header)
	h.Retain = b&RetainPos == 1
	h.Qos = QoSLevel((b & QosPos) >> 1)
	h.Dup = (b&DupPos)>>3 == 1
	h.MsgType = MessageType((b & MsgTypePos) >> 4)
	h.RemainLength, err = d.readRemainLength()
	return h, err
}

func (d *decoder) readRemainLength() (int, error) {
	multi, value := 1, 0
	for {
		b, err := d.ReadByte()
		if err != nil {
			return 0, err
		}
		value += (int(b) & 127) * multi
		multi *= 128

		// 最高位即 continuation bit 为0，表明是最后一字节
		if (b & 0xff) == 0 {
			break
		}
	}

	return value, nil
}

func (d *decoder) readFixedBytes() ([]byte, error) {
	fixedLen, err := d.readUint16()
	if err != nil {
		return nil, err
	}

	buf := make([]byte, fixedLen)
	_, err = io.ReadFull(d, buf)
	return buf, err
}

func (d *decoder) readUint16() (uint16, error) {
	bs := make([]byte, 2)
	if _, err := d.Read(bs); err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(bs), nil
}

func (d *decoder) readOneByte() (byte, error) {
	return d.ReadByte()
}

type encoder struct {
	*bufio.Writer
}

func (e *encoder) writeHeader(h *Header) error {
	b := h.MsgType.Byte()<<4 |
		boolToByte(h.Dup)<<3 |
		h.Qos.Byte()<<1 |
		boolToByte(h.Retain)<<1
	if err := e.WriteByte(b); err != nil {
		return err
	}
	return e.writeRemainLength(h.RemainLength)
}

func (e *encoder) writeRemainLength(n int) error {
	for {
		b := n % 128
		n /= 128
		if n > 0 {
			b |= 0x80
		}
		if err := e.WriteByte(byte(b)); err != nil {
			return err
		}
		if n == 0 {
			break
		}
	}
	return nil
}

func boolToByte(b bool) byte {
	if b {
		return 1
	}
	return 0
}
