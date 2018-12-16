package relay

import (
	"github.com/0990/gorpc"
	"github.com/0990/gorpc/codec"
)

type PassthroughContext struct {
	Int64      int64
	Int64Slice []int64
	Str        string
}

func ResolveInboundEvent(inputEvent gorpc.Event) (outputEvent gorpc.Event, handled bool, err error) {
	switch relayMsg := inputEvent.Message().(type) {
	case *RelayACK:
		ev := &RecvMsgEvent{
			Ses: inputEvent.Session(),
			ack: relayMsg,
		}
		if relayMsg.MsgID != 0 {
			ev.Msg, _, err = codec.DecodeMessage(int(relayMsg.MsgID), relayMsg.Msg)
			if err != nil {
				return
			}
		}

		return ev, true, nil
	}
	return inputEvent, false, nil
}

func ResolveOutboundEvent(inputEvent gorpc.Event) (handled bool, err error) {
	switch relayMsg := inputEvent.Message().(type) {
	case *RelayACK:
		if relayMsg.MsgID != 0 {
			_, _, err = codec.DecodeMessage(int(relayMsg.MsgID), relayMsg.Msg)
			if err != nil {
				return
			}
		}
		return true, nil
	}
	return
}
