package rpc

import (
	"errors"
	"github.com/0990/gorpc"
	"github.com/0990/gorpc/relay"
)

var (
	ErrInvalidPeerSession = errors.New("rpc: Invalid peer type, require cellnet.RPCSessionGetter or cellnet.Session")
	ErrEmptySession       = errors.New("rpc: Empty session")
)

type RPCSessionGetter interface {
	RPCSession() gorpc.Session
}

func getPeerSession(ud interface{}) (ses gorpc.Session, err error) {
	if ud == nil {
		return nil, relay.ErrInvalidPeerSession
	}

	switch i := ud.(type) {
	case RPCSessionGetter:
		ses = i.RPCSession()
	case gorpc.Session:
		ses = i
	case gorpc.TCPConnector:
		ses = i.Session()
	default:
		err = ErrInvalidPeerSession
		return
	}

	if ses == nil {
		return nil, ErrEmptySession
	}
	return
}
