package tcp

import (
	"github.com/0990/gorpc"
	"github.com/0990/gorpc/peer"
	"github.com/0990/gorpc/util"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type tcpSession struct {
	peer.CoreContextSet
	peer.CoreSessionIdentify
	*peer.CoreProcBundle
	pInterface gorpc.Peer

	conn net.Conn

	exitSync sync.WaitGroup

	sendQueue    *gorpc.Pipe
	cleanupGurad sync.Mutex
	endNotify    func()
	closing      int64
}

func (self *tcpSession) Peer() gorpc.Peer {
	return self.pInterface
}

func (self *tcpSession) Raw() interface{} {
	return self.conn
}

func (self *tcpSession) Close() {
	closing := atomic.SwapInt64(&self.closing, 1)

	if closing != 0 {
		return
	}

	if self.conn != nil {
		con := self.conn.(*net.TCPConn)
		con.CloseRead()
		con.SetReadDeadline(time.Now())
	}
}

func (self *tcpSession) Send(msg interface{}) {
	if msg == nil {
		return
	}

	if self.IsManualClosed() {
		return
	}

	self.sendQueue.Add(msg)
}

func (self *tcpSession) IsManualClosed() bool {
	return atomic.LoadInt64(&self.closing) != 0
}

func (self *tcpSession) protectedReadMessage() (msg interface{}, err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("IO panic: %s", err)
			self.conn.Close()
		}
	}()

	msg, err = self.ReadMessage(self)
	return
}

func (self *tcpSession) recvLoop() {
	var capturePanic bool
	if i, ok := self.Peer().(gorpc.PeerCaptureIOPanic); ok {
		capturePanic = i.CaptureIOPanic()
	}

	for self.conn != nil {
		var msg interface{}
		var err error
		if capturePanic {
			msg, err = self.protectedReadMessage()
		} else {
			msg, err = self.ReadMessage(self)
		}

		if err != nil {
			if !util.IsEOFOrNetReadError(err) {
				log.Errorf("session closed,sesid:%d,err:%s", self.ID(), err)
			}
			self.sendQueue.Add(nil)

			closedMsg := &gorpc.SessionClosed{}

			if self.IsManualClosed() {
				closedMsg.Reason = gorpc.CloseReason_Manual
			}

			self.ProcEvent(&gorpc.RecvMsgEvent{
				Ses: self,
				Msg: closedMsg,
			})
			break
		}
		self.ProcEvent(&gorpc.RecvMsgEvent{Ses: self, Msg: msg})
	}

	self.exitSync.Done()
}

func (self *tcpSession) sendLoop() {
	var writeList []interface{}
	for {
		writeList = writeList[0:0]
		exit := self.sendQueue.Pick(&writeList)

		for _, msg := range writeList {
			self.SendMessage(&gorpc.SendMsgEvent{
				Ses: self,
				Msg: msg,
			})
		}
		if exit {
			break
		}
	}

	self.conn.Close()
	self.exitSync.Done()
}

func (self *tcpSession) Start() {
	atomic.StoreInt64(&self.closing, 0)

	self.sendQueue.Reset()
	self.exitSync.Add(2)

	self.Peer().(peer.SessionManager).Add(self)

	go func() {
		self.exitSync.Wait()
		self.Peer().(peer.SessionManager).Remove(self)

		if self.endNotify != nil {
			self.endNotify()
		}
	}()

	go self.recvLoop()
	go self.sendLoop()

}

func newSession(conn net.Conn, p gorpc.Peer, endNotify func()) *tcpSession {
	self := &tcpSession{
		conn:       conn,
		endNotify:  endNotify,
		sendQueue:  gorpc.NewPipe(),
		pInterface: p,
		CoreProcBundle: p.(interface {
			GetBundle() *peer.CoreProcBundle
		}).GetBundle(),
	}
	return self
}
