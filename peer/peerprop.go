package peer

import "github.com/0990/gorpc"

type CorePeerProperty struct {
	name  string
	queue gorpc.EventQueue
	addr  string
}

func (self *CorePeerProperty) Name() string {
	return self.name
}

func (self *CorePeerProperty) Queue() gorpc.EventQueue {
	return self.queue
}

func (self *CorePeerProperty) Address() string {
	return self.addr
}

func (self *CorePeerProperty) SetName(v string) {
	self.name = v
}

func (self *CorePeerProperty) SetQueue(v gorpc.EventQueue) {
	self.queue = v
}

func (self *CorePeerProperty) SetAddress(v string) {
	self.addr = v
}
