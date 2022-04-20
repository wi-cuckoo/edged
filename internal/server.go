package internal

import (
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/wi-cuckoo/edged"
)

func Setup(c *cli.Context) error {
	setLogger(c)

	ln, err := net.Listen("tcp", c.String(edged.TCPAddrFlag.Name))
	if err != nil {
		return err
	}
	logrus.Infof("listen on tcp address %s", ln.Addr().String())

	e := &Edged{
		ln:  ln,
		cfg: c,
		ps: &pubsub{
			topics: make(map[string]*topic),
		},
		quit: make(chan struct{}),
	}
	if err := e.serve(); err != nil {
		return err
	}
	e.waitClose()

	return nil
}

type Edged struct {
	cfg  *cli.Context
	ln   net.Listener
	ps   *pubsub
	quit chan struct{}
}

func (e *Edged) serve() error {
	for {
		conn, err := e.ln.Accept()
		if err != nil {
			return err
		}
		logrus.Infof("accept publisher conn %s->%s", conn.RemoteAddr(), conn.LocalAddr())
		edgedConn := NewEdgeConn(conn, e.ps)
		go edgedConn.handleConn()
	}
}

// WaitClose listen to sys singal, then do something befor exit really
func (e *Edged) waitClose() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	logrus.Info("receive stop signal from system, exiting......")

	close(e.quit)
	done := make(chan bool)
	go func() {
		e.ln.Close()
		done <- true
	}()

	select {
	case <-sig:
		logrus.Warn("receive close signal second times")
	case <-time.After(time.Second * 3):
		logrus.Warn("3s elapsed, shutdown hardly")
	case <-done:
		logrus.Debug("goodbye, exited successfully")
	}
}

func setLogger(c *cli.Context) {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)
	if c.Bool(edged.DebugFlag.Name) {
		logrus.SetLevel(logrus.DebugLevel)
	}
}
