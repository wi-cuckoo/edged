package internal

import (
	"net"

	"github.com/sirupsen/logrus"
)

type Sub struct {
	ln net.Listener
}

func NewSub(ln net.Listener) *Sub {
	return &Sub{ln}
}

func (s *Sub) Serve() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			return
		}
		logrus.Infof("accept subscriber conn %s->%s", conn.RemoteAddr(), conn.LocalAddr())
		go s.handleConn(conn)
	}
}

func (s *Sub) Close() {
	s.ln.Close()
}

func (s *Sub) handleConn(conn net.Conn) {

}
