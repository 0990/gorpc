package tcp

import (
	"github.com/0990/gorpc"
	"github.com/0990/gorpc/proc"
	"github.com/davyxu/cellnet/proc/tcp"
)

func init() {
	proc.RegisterProcessor("tcp.ltv", func(bundle proc.ProcessorBundle, userCallback gorpc.EventCallback) {
		bundle.SetTransmitter(new(tcp.TCPMessageTransmitter))
		bundle.SetHooker(new(tcp.MsgHooker))
		bundle.SetCallback(proc.NewQueuedEventCallback(userCallback))
	})
}
