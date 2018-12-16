package tcp

import (
	"github.com/0990/gorpc"
	"github.com/davyxu/cellnet/msglog"
	"github.com/davyxu/cellnet/relay"
	"github.com/davyxu/cellnet/rpc"
)

type MsgHooker struct {
}

func (m MsgHooker) OnInboundEvent(inputEvent gorpc.Event) (output gorpc.Event) {
	var handled bool
	var err error

	inputEvent, handled, err = rpc.ResolveInboundEvent(inputEvent)
	if err != nil {
		log.Errorf("rpc.ResolveInboundEvent", err)
		return
	}

	if !handled {
		inputEvent, handled, err = relay.ResoleveInboundEvent(inputEvent)
		if err != nil {
			log.Errorln("relay.ResolveInboundEvent:", err)
			return
		}

		if !handled {
			msglog.WriteRecvLogger(log, "tcp", inputEvent.Session(), inputEvent.Messag())
		}
	}
}
