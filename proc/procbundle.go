package proc

import "github.com/0990/gorpc"

type ProcessorBundle interface {
	SetTransmitter(v gorpc.MessageTransmitter)

	SetHooker(v gorpc.EventHooker)

	SetCallback(v gorpc.EventCallback)
}

//execute callback in session's peer's queue
func NewQueuedEventCallback(callback gorpc.EventCallback) gorpc.EventCallback {
	return func(ev gorpc.Event) {
		if callback != nil {
			gorpc.SessionQueuedCall(ev.Session(), func() {
				callback(ev)
			})
		}
	}
}

type MultiHooker []gorpc.EventHooker

func (mh MultiHooker) OnInboundEvent(input gorpc.Event) (output gorpc.Event) {
	for _, h := range mh {
		input = h.OnInboundEvent(input)
		if input == nil {
			break
		}
	}

	return input
}

func (mh MultiHooker) OnOutboundEvent(input gorpc.Event) (output gorpc.Event) {
	for _, h := range mh {
		input = h.OnOutboundEvent(input)
		if input == nil {
			break
		}
	}
	return input
}

func NewMultiHooker(h ...gorpc.EventHooker) gorpc.EventHooker {
	return MultiHooker(h)
}
