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
	"github.com/wi-cuckoo/edged/internal"
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
		pub:  internal.NewPub(pubLn),
		sub:  internal.NewSub(subLn),
		quit: make(chan struct{}),
	}
	if err := e.Start(); err != nil {
		return err
	}
	e.WaitClose()

	return nil
}

type Edged struct {
	pub  *internal.Pub
	sub  *internal.Sub
	quit chan struct{}
}

func (e *Edged) Start() error {
	go e.pub.Serve()
	go e.sub.Serve()

	return nil
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
