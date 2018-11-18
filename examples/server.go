package main

import (
	"fmt"
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/proc"
	"github.com/davyxu/cellnet/rpc"
	"github.com/davyxu/cellnet/tests"
)

func server() {
	queue := cellnet.NewEventQueue()
	peerIns := peer.NewGenericPeer("tcp.Acceptor", "server", peerAddress, queue)

	proc.BindProcessorHandler(peerIns, "tcp.ltv", func(ev cellnet.Event) {
		switch msg := ev.Message().(type) {
		case *cellnet.SessionAccepted:
			fmt.Println("server accepted")
		case *tests.TestEchoACK:
			fmt.Println("server rece %+v\n", msg)
			ack := &tests.TestEchoACK{
				Msg:   msg.Msg,
				Value: msg.Value,
			}
			if rpcevent, ok := ev.(*rpc.RecvMsgEvent); ok {
				rpcevent.Reply(ack)
			} else {
				ev.Session().Send(ack)
			}
		case *cellnet.SessionClosed:
			fmt.Println("session closed:", ev.Session().ID())
		}
	})
	peerIns.Start()
	queue.StartLoop()
}
