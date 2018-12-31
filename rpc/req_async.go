package rpc

import (
	"github.com/0990/gorpc"
	"github.com/sirupsen/logrus"
	"time"
)

func Call(sesOrPeer interface{}, reqMsg interface{}, timeout time.Duration, callback func(raw interface{})) {
	ses, err := getPeerSession(sesOrPeer)
	if err != nil {
		logrus.Errorln("Remote call failed", err)
		gorpc.SessionQueuedCall(ses, func() {
			callback(err)
		})
	}

	req := createRequest(func(raw interface{}) {
		gorpc.SessionQueuedCall(ses, func() {
			callback(raw)
		})
	})

	req.Send(ses, reqMsg)

	time.AfterFunc(timeout, func() {
		if getRequest(req.id) != nil {
			gorpc.SessionQueuedCall(ses, func() {
				callback(ErrTimeout)
			})
		}
	})
}
