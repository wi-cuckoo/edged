package internal

import (
	"net"

	"github.com/sirupsen/logrus"
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

func (p *Pub) handleConn(conn net.Conn) {

}
