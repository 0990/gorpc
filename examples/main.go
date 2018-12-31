package main

import (
	"fmt"
	"github.com/0990/gorpc"
	"github.com/0990/gorpc/codec"
	"github.com/0990/gorpc/peer"
	_ "github.com/0990/gorpc/peer/tcp" // 注册TCP Peer
	"github.com/0990/gorpc/proc"
	_ "github.com/0990/gorpc/proc/tcp" // 注册TCP Processor
	"github.com/0990/gorpc/rpc"
	"github.com/0990/gorpc/util"
	"reflect"
	"time"
)

const peerAddress = "127.0.0.1:17701"

func main() {
	server()
	//send
	//clientAsyncCallback()
	//request
	//clientAsyncRPC()
	//call
	clientSyncRPC()
}

func server() {
	queue := gorpc.NewEventQueue()
	peerIns := peer.NewGenericPeer("tcp.Acceptor", "server", peerAddress, queue)

	proc.BindProcessorHandler(peerIns, "tcp.ltv", func(ev gorpc.Event) {
		switch msg := ev.Message().(type) {
		case *gorpc.SessionAccepted:
			fmt.Println("server accepted")
		case *TestEchoAck:
			fmt.Printf("server recv %+v\n", msg)
			ack := &TestEchoAck{
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
		default:
			fmt.Printf("default server recv %+v\n", msg)
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
			ev.Session().Send(&TestEchoAck{
				Msg:   "hello",
				Value: 1234,
			})
		case *TestEchoAck:
			fmt.Printf("clientAsyncCallback recv %+v\n", msg)
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

func clientAsyncRPC() {
	done := make(chan struct{})
	queue := gorpc.NewEventQueue()

	p := peer.NewGenericPeer("tcp.Connector", "async rpc", peerAddress, queue)

	rv := proc.NewSyncReceiver()

	proc.BindProcessorHandler(p, "tcp.ltv", rv.EventCallback())

	p.Start()
	queue.StartLoop()
	rv.WaitMessage("gorpc.SessionConnected")

	rpc.Call(p, &TestEchoAck{
		Msg:   "hello",
		Value: 12334,
	}, 10*time.Second, func(raw interface{}) {
		switch result := raw.(type) {
		case error:
			fmt.Println(result)
		default:
			fmt.Println(result)
			done <- struct{}{}
		}
	})
	<-done
}

type TestEchoAck struct {
	Msg   string
	Value int32
}

func (self *TestEchoAck) String() string {
	return fmt.Sprintf("%+v", *self)
}

func init() {
	gorpc.RegisterMessageMeta(&gorpc.MessageMeta{
		Codec: codec.MustGetCodec("binary"),
		Type:  reflect.TypeOf((*TestEchoAck)(nil)).Elem(),
		ID:    int(util.StringHash("main.TestEchoAck")),
	})
}

func clientSyncRPC() {
	queue := gorpc.NewEventQueue()
	p := peer.NewGenericPeer("tcp.Connector", "async rpc", peerAddress, queue)

	rv := proc.NewSyncReceiver()
	proc.BindProcessorHandler(p, "tcp.ltv", rv.EventCallback())
	p.Start()

	queue.StartLoop()

	rv.WaitMessage("gorpc.SessionConnected")
	rpc.CallSync(p, &TestEchoAck{
		Msg:   "rpc sync",
		Value: 1234,
	}, time.Second)
}
