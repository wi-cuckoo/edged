package protocol

import (
	"bufio"
	"io"
)

// ReadPacket ...
func ReadPacket(r io.Reader) error {
	dec := &decoder{bufio.NewReader(r)}
	header, err := dec.readHeader()
	if err != nil {
		return err
	}

	if header.MsgType == CONNECT {
	}
	return nil
}
