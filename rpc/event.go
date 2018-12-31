package rpc

import (
	"github.com/0990/gorpc"
	"github.com/0990/gorpc/codec"
	"github.com/sirupsen/logrus"
)

type RecvMsgEvent struct {
	ses    gorpc.Session
	Msg    interface{}
	callid int64
}

func (self *RecvMsgEvent) Session() gorpc.Session {
	return self.ses
}

func (self *RecvMsgEvent) Message() interface{} {
	return self.Msg
}

func (self *RecvMsgEvent) Queue() gorpc.EventQueue {
	return self.ses.Peer().(interface {
		Queue() gorpc.EventQueue
	}).Queue()
}

func (self *RecvMsgEvent) Reply(msg interface{}) {
	data, meta, err := codec.EncodeMessage(msg)
	if err != nil {
		logrus.Errorf("rpc reply message encode error:%s", err)
		return
	}
	self.ses.Send(&RemoteCallACK{
		MsgID:  uint32(meta.ID),
		Data:   data,
		CallID: self.callid,
	})
}
