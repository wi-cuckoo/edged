package internal

import (
	"net"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/wi-cuckoo/edged/protocol"
)

type EdgedConn struct {
	net.Conn

	authrized bool
	ps        *pubsub
	topic     *topic
	quit      chan struct{}

	wg      sync.WaitGroup
	in, out chan protocol.Packet
}

func NewEdgeConn(conn net.Conn, ps *pubsub) *EdgedConn {
	return &EdgedConn{
		Conn:      conn,
		ps:        ps,
		authrized: false,
		quit:      make(chan struct{}),
		in:        make(chan protocol.Packet, 16),
		out:       make(chan protocol.Packet, 16),
	}
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
}

func (c *EdgedConn) processIncoming() {
	defer c.wg.Done()
	for {
		packet, err := protocol.ReadPacket(c)
		if err != nil {
			logrus.Errorf("read packet err: %s", err.Error())
			break
		}
		switch packet.MessageType() {
		case protocol.CONNECT:
			// return CONNACK package
			c.out <- c.handleConnect(packet)
		case protocol.DISCONNECT:
			return
		case protocol.PINGREQ:
			// return PINGRESP packet
		case protocol.PUBLISH:
			// return PUBACK packet
		case protocol.SUBSCRIBE:
			// return SUBACK packet
		case protocol.UNSUBSCRIBE:
			// return UNSUBACK packet
		}
	}
}

func (c *EdgedConn) processOutgoing() {
	defer c.wg.Done()
	defer c.Close()

	for {
		select {
		case packet := <-c.out:
			// write packet to client
			if err := protocol.WritePacket(c, packet); err != nil {
				logrus.Errorf("write packet fail: %s", err.Error())
				return
			}
		case <-c.quit:
			return
		}
	}
}

func (c *EdgedConn) handleConnect(p protocol.Packet) *protocol.ConnackPacket {
	cp := p.(*protocol.ConnectPacket)

	if len(cp.WillTopic) == 0 {
		return nil
	}
	c.topic = c.ps.createTopic(string(cp.WillTopic))
	c.authrized = true // todo check passws&username

	ca := &protocol.ConnackPacket{
		Header:  cp.Header,
		RetCode: protocol.Accepted,
	}

	return ca
}
