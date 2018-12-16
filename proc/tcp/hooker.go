package tcp

import (
	"fmt"
	"github.com/0990/gorpc"
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
		//inputEvent, handled, err = relay.ResoleveInboundEvent(inputEvent)
		//if err != nil {
		//	log.Errorln("relay.ResolveInboundEvent:", err)
		//	return
		//}

		//if !handled {
		//	msglog.WriteRecvLogger(log, "tcp", inputEvent.Session(), inputEvent.Message())
		//}
		fmt.Println("msg unhandled")
	}
	return
}

func (self MsgHooker) OnOutboundEvent(inputEvent gorpc.Event) (outputEvent gorpc.Event) {
	handled, err := rpc.ResolveOutboundEvent(inputEvent)

	if err != nil {
		log.Errorln("repc.ResolveOutboundEvent:", err)
		return nil
	}

	if !handled {
		fmt.Println("msg unhandled")
	}
	return
}
