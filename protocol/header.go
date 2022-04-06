package protocol

const (
	RetainPos  = 0x01
	QosPos     = 0x06
	DupPos     = 0x08
	MsgTypePos = 0xf0
)

// QoSLevel Position: byte 1, bits 2-1
type QoSLevel byte

const (
	QoSLevel0 = iota // At most once, Fire and Forget
	QoSLevel1        // At least once, Acknowledged delivery
	QoSLevel2        // Exactly once, Assured delivery
)

// MessageType Position: byte 1, bits 7-4
type MessageType byte

const (
	Reserved    MessageType = iota
	CONNECT                 //Client request to connect to Server
	CONNACK                 // Connect Acknowledgment
	PUBLISH                 // Publish message
	PUBACK                  // Publish Acknowledgment
	PUBREC                  // Publish Received (assured delivery part 1)
	PUBREL                  // Publish Release (assured delivery part 2)
	PUBCOMP                 // Publish Complete (assured delivery part 3)
	SUBSCRIBE               // Client Subscribe request
	SUBACK                  // Subscribe Acknowledgment
	UNSUBSCRIBE             // Client Unsubscribe request
	UNSUBACK                // Unsubscribe Acknowledgment
	PINGREQ                 // PING Request
	PINGRESP                // PING Response
	DISCONNECT              // Client is Disconnecting
)

// Header fixed byte 1
type Header struct {
	MsgType      MessageType
	Qos          QoSLevel
	Retain, Dup  bool
	RemainLength int
}
