package relay

import "github.com/0990/gorpc"

type RecvMsgEvent struct {
	Ses gorpc.Session
	ack *RelayACK
	Msg interface{}
}

func (self *RecvMsgEvent) PassThroughAsInt64() int64 {
	if self.ack == nil {
		return 0
	}
	return self.ack.Int64
}

func (self *RecvMsgEvent) PassThroughAsInt64Slice() []int64 {
	if self.ack == nil {
		return nil
	}
	return self.ack.Int64Slice
}

func (self *RecvMsgEvent) PassThroughAsString() string {
	if self.ack == nil {
		return ""
	}
	return self.ack.Str
}

func (self *RecvMsgEvent) Session() gorpc.Session {
	return self.Ses
}

func (self *RecvMsgEvent) Message() interface{} {
	return self.Msg
}

func (self *RecvMsgEvent) Reply(msg interface{}) {
	Relay(self.Ses, msg, self.ack.Int64, self.ack.Int64Slice, self.ack.Str)
}
