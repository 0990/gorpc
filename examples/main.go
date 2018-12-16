package main

import (
	"fmt"
	"github.com/0990/gorpc"
	"github.com/0990/gorpc/peer"
	_ "github.com/0990/gorpc/peer/tcp" // 注册TCP Peer
	"github.com/0990/gorpc/proc"
	_ "github.com/0990/gorpc/proc/tcp" // 注册TCP Processor
	"github.com/0990/gorpc/rpc"
	"github.com/0990/gorpc/tests"
)

const peerAddress = "127.0.0.1:17701"

func main() {
	server()

	clientAsyncCallback()
}

func server() {
	queue := gorpc.NewEventQueue()
	peerIns := peer.NewGenericPeer("tcp.Acceptor", "server", peerAddress, queue)

	proc.BindProcessorHandler(peerIns, "tcp.ltv", func(ev gorpc.Event) {
		switch msg := ev.Message().(type) {
		case *gorpc.SessionAccepted:
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
		case *gorpc.SessionClosed:
			fmt.Println("session closed:", ev.Session().ID())
		}
	})
	peerIns.Start()
	queue.StartLoop()
}

func clientAsyncCallback() {
	done := make(chan struct{})
	queue := gorpc.NewEventQueue()

	p := peer.NewGenericPeer("tcp.Connector", "clientAsyncCallback", peerAddress, queue)

	proc.BindProcessorHandler(p, "tcp.ltv", func(ev gorpc.Event) {
		switch msg := ev.Message().(type) {
		case *gorpc.SessionConnected:
			fmt.Println("clientAsyncCallback connected")
			ev.Session().Send(&tests.TestEchoACK{
				Msg:   "hello",
				Value: 1234,
			})
		case *tests.TestEchoACK:
			fmt.Println("clientAsyncCallback recv %+v\n", msg)
			done <- struct{}{}
		case *gorpc.SessionClosed:
			fmt.Println("clientAsyncCallback closed")
		default:
			fmt.Println(msg)
		}
	})
	p.Start()
	queue.StartLoop()
	<-done
}
