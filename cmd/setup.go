package main

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

func setup(c *cli.Context) error {
	setLogger(c)

	pubLn, err := net.Listen("tcp", c.String(edged.PubAddrFlag.Name))
	if err != nil {
		return err
	}
	subLn, err := net.Listen("tcp", c.String(edged.SubAddrFlag.Name))
	if err != nil {
		return err
	}

	e := &Edged{
		pub:  pubLn,
		sub:  subLn,
		quit: make(chan struct{}),
	}
	go e.servePub()
	go e.serveSub()
	e.WaitClose()

	return nil
}

type Edged struct {
	pub, sub net.Listener
	quit     chan struct{}
}

func (e *Edged) servePub() {
	for {
		conn, err := e.pub.Accept()
		if err != nil {
			return
		}
		logrus.Infof("accept conn %s->%s", conn.RemoteAddr(), conn.LocalAddr())
	}
}

func (e *Edged) serveSub() {
	for {
		conn, err := e.pub.Accept()
		if err != nil {
			return
		}
		logrus.Infof("accept conn %s->%s", conn.RemoteAddr(), conn.LocalAddr())
	}
}

// WaitClose listen to sys singal, then do something befor exit really
func (e *Edged) WaitClose() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	logrus.Info("receive stop signal from system, exiting......")

	close(e.quit)
	done := make(chan bool)
	go func() {
		e.pub.Close()
		e.sub.Close()
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
