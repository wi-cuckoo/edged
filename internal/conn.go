package internal

import (
	"errors"
	"fmt"
	"net"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/wi-cuckoo/edged/protocol"
)

type EdgedConn struct {
	net.Conn

	authrized bool
	topic     *Topic
	quit      chan struct{}

	wg      sync.WaitGroup
	in, out chan protocol.Packet
}

func (c *EdgedConn) handleConn() {
	defer c.close()

	c.wg.Add(2)
	go c.processIncoming()
	go c.processOutgoing()
	c.wg.Wait()

}

func (c *EdgedConn) close() {
	close(c.in)
	close(c.out)
	c.Close()
}

func (c *EdgedConn) processIncoming() {
	defer c.wg.Done()
	for {
		packet, err := protocol.ReadPacket(c)
		if err != nil {
			logrus.Errorf("read packet err: %s", err.Error())
			break
		}
		c.in <- packet
	}
}

func (c *EdgedConn) processOutgoing() {
	defer c.wg.Done()
	for {
		select {
		case packet := <-c.out:
			// write packet to client
			fmt.Println(packet)
		case <-c.quit:
			return
		}
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
