package tcp

import (
	"github.com/0990/gorpc"
	"github.com/0990/gorpc/proc"
)

func init() {
	proc.RegisterProcessor("tcp.ltv", func(bundle proc.ProcessorBundle, userCallback gorpc.EventCallback) {
		bundle.SetTransmitter(new(TCPMessageTransmitter))
		bundle.SetHooker(new(MsgHooker))
		bundle.SetCallback(proc.NewQueuedEventCallback(userCallback))
	})
}
