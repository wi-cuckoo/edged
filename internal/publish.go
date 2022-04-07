package internal

import (
	"errors"
	"net"

	"github.com/sirupsen/logrus"
	"github.com/wi-cuckoo/edged/protocol"
)

type Pub struct {
	ln net.Listener
}

func NewPub(ln net.Listener) *Pub {
	return &Pub{ln}
}

func (p *Pub) Close() {
	p.ln.Close()
}

func (p *Pub) Serve() {
	for {
		conn, err := p.ln.Accept()
		if err != nil {
			return
		}
		logrus.Infof("accept publisher conn %s->%s", conn.RemoteAddr(), conn.LocalAddr())
		go p.handleConn(conn)
	}
}

type PubConn struct {
	authrized bool
	topic     *Topic
	conn      net.Conn
}

func (p *Pub) handleConn(conn net.Conn) {
	defer conn.Close()

	var authrized bool
	for {
		packet, err := protocol.ReadPacket(conn)
		if err != nil {
			logrus.Errorf("read packet err: %s", err.Error())
			return
		}

		if !authrized {
			if err := authenticate(packet); err != nil {
				logrus.Errorf("auth fail err: %s", err.Error())
				return
			}
			authrized = true
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
