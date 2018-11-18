package main

import (
	"fmt"
	"github.com/0990/gorpc"
	"github.com/0990/gorpc/peer"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proc"
	"github.com/davyxu/cellnet/tests"
)

func clientAsyncCallback() {
	done := make(chan struct{})
	queue := gorpc.NewEventQueue()

	p := peer.NewGenericPeer("tcp.Connector", "clientAsyncCallback", peerAddress, queue)

	proc.BindProcessorHandler(p, "tcp.ltv", func(ev cellnet.Event) {
		switch msg := ev.Message().(type) {
		case *cellnet.SessionConnected:
			fmt.Println("clientAsyncCallback connected")
			ev.Session().Send(&tests.TestEchoACK{
				Msg:   "hello",
				Value: 1234,
			})
		case *tests.TestEchoACK:
			fmt.Println("clientAsyncCallback recv %+v\n", msg)
			done <- struct{}{}
		case *cellnet.SessionClosed:
			fmt.Println("clientAsyncCallback closed")
		}
	})
	p.Start()
	queue.StartLoop()
	<-done
}
