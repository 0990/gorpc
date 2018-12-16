package gorpc

import "github.com/davyxu/cellnet"

type Session interface {
	Raw() interface{}

	Peer() Peer
	Send(msg interface{})
	Close()
	ID() int64
}

type RawPacket struct {
	MsgData []byte
	MsgID   int
}

func (self *RawPacket) Messge() interface{} {
	meta := cellnet.MessageMetaByID(self.MsgID)

	if meta == nil {
		return struct{}{}
	}

	msg := meta.NewType()
	err := meta.Codec.Decode(self.MsgData, msg)

	if err != nil {
		return struct{}{}
	}

	return msg
}
