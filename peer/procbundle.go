package peer

import (
	"errors"
	"github.com/0990/gorpc"
)

type MessagePoster interface {
	ProcEvent(ev gorpc.Event)
}

type CoreProcBundle struct {
	transmit gorpc.MessageTransmitter
	hooker   gorpc.EventHooker
	callback gorpc.EventCallback
}

func (self *CoreProcBundle) GetBundle() *CoreProcBundle {
	return self
}

func (self *CoreProcBundle) SetTransmitter(v gorpc.MessageTransmitter) {
	self.transmit = v
}

func (self *CoreProcBundle) SetHooker(v gorpc.EventHooker) {
	self.hooker = v
}

func (self *CoreProcBundle) SetCallback(v gorpc.EventCallback) {
	self.callback = v
}

var notHandled = errors.New("Processor: Transimitter nil")

func (self *CoreProcBundle) ReadMessage(ses gorpc.Session) (msg interface{}, err error) {

	if self.transmit != nil {
		return self.transmit.OnRecvMessage(ses)
	}

	return nil, notHandled
}

func (self *CoreProcBundle) SendMessage(ev gorpc.Event) {

	if self.hooker != nil {
		ev = self.hooker.OnOutboundEvent(ev)
	}

	if self.transmit != nil && ev != nil {
		self.transmit.OnSendMessage(ev.Session(), ev.Message())
	}
}

func (self *CoreProcBundle) ProcEvent(ev gorpc.Event) {

	if self.hooker != nil {
		ev = self.hooker.OnInboundEvent(ev)
	}

	if self.callback != nil && ev != nil {
		self.callback(ev)
	}
}
