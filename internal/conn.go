package internal

import (
	"errors"
	"net"

	"github.com/sirupsen/logrus"
	"github.com/wi-cuckoo/edged/protocol"
)

type EdgedConn struct {
	net.Conn

	authrized bool
	topic     *Topic
	quit      chan struct{}

	in, out chan protocol.Packet
}

func (c *EdgedConn) handleConn() {
	defer c.Close()

	for {
		packet, err := protocol.ReadPacket(c)
		if err != nil {
			logrus.Errorf("read packet err: %s", err.Error())
			return
		}

		if !c.authrized {
			if err := authenticate(packet); err != nil {
				logrus.Errorf("auth fail err: %s", err.Error())
				return
			}
			c.authrized = true
			// return ack
			continue
		}

		// todo handle publish
	}
}

func authenticate(p protocol.Packet) error {
	_, ok := p.(*protocol.ConnectPacket)
	if !ok {
		return errors.New("not CONNECT packet")
	}
	// todo check username & password

	return nil
}
