package rpc

import (
	"github.com/0990/gorpc"
	"github.com/0990/gorpc/codec"
)

type RemoteCallMsg interface {
	GetMsgID() uint16
	GetMsgData() []byte
	GetCallID() int64
}

func ResolveInboundEvent(inputEvent gorpc.Event) (outputEvent gorpc.Event, handled bool, err error) {
	if _, ok := inputEvent.(*RecvMsgEvent); ok {
		return inputEvent, false, nil
	}
	rpcMsg, ok := inputEvent.Message().(RemoteCallMsg)
	if !ok {
		return inputEvent, false, nil
	}

	userMsg, _, err := codec.DecodeMessage(int(rpcMsg.GetMsgID()), rpcMsg.GetMsgData())
	if err != nil {
		return inputEvent, false, err
	}
	if log.IsDebugEnabled() {
		peerInfo := inputEvent.Session().Peer().(gorpc.PeerProperty)

		log.Debugf("#rpc.recv(%s)@%d len: %d %s | %s",
			peerInfo.Name(),
			inputEvent.Session().ID(),
			gorpc.MessageSize(userMsg),
			gorpc.MessageToName(userMsg),
			gorpc.MessageToString(userMsg))
	}

	switch inputEvent.Message().(type) {
	case *RemoteCallREQ:
		return &RecvMsgEvent{
			inputEvent.Session(),
			userMsg,
			rpcMsg.GetCallID(),
		}, true, nil
	case *RemoteCallACK:
		request := getRequest(rpcMsg.GetCallID())
		if request != nil {
			request.RecvFeedback(userMsg)
		}
		return inputEvent, true, nil

	}
	return inputEvent, false, nil
}

func ResolveOutboundEvent(inputEvent gorpc.Event) (handled bool, err error) {
	rpcMsg, ok := inputEvent.Message().(RemoteCallMsg)
	if !ok {
		return false, nil
	}
	userMsg, _, err := codec.DecodeMessage(int(rpcMsg.GetMsgID()), rpcMsg.GetMsgData())

	if err != nil {
		return false, err
	}

	if log.IsDebugEnabled() {

		peerInfo := inputEvent.Session().Peer().(gorpc.PeerProperty)

		log.Debugf("#rpc.send(%s)@%d len: %d %s | %s",
			peerInfo.Name(),
			inputEvent.Session().ID(),
			gorpc.MessageSize(userMsg),
			gorpc.MessageToName(userMsg),
			gorpc.MessageToString(userMsg))
	}

	// 避免后续环节处理

	return true, nil

}
