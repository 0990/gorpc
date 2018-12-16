package tcp

import (
	"fmt"
	"github.com/0990/gorpc"
	"github.com/0990/gorpc/relay"
	"github.com/0990/gorpc/rpc"
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
		inputEvent, handled, err = relay.ResolveInboundEvent(inputEvent)
		if err != nil {
			log.Errorln("relay.ResolveInboundEvent:", err)
			return
		}
		if !handled {
			fmt.Println("msg unhandled")
		}

		//if !handled {
		//	msglog.WriteRecvLogger(log, "tcp", inputEvent.Session(), inputEvent.Message())
		//}

	}
	return inputEvent
}

func (self MsgHooker) OnOutboundEvent(inputEvent gorpc.Event) (outputEvent gorpc.Event) {
	handled, err := rpc.ResolveOutboundEvent(inputEvent)

	if err != nil {
		log.Errorln("repc.ResolveOutboundEvent:", err)
		return nil
	}

	if !handled {
		handled, err = relay.ResolveOutboundEvent(inputEvent)
		if err != nil {
			log.Errorln("relay.ResolveOutboundEvent:", err)
		}
		if !handled {
			fmt.Println("msg unhandled")
		}
	}
	return inputEvent
}
