package tcp

import (
	"github.com/0990/gorpc"
	"github.com/0990/gorpc/peer"
	"github.com/0990/gorpc/util"
	"github.com/sirupsen/logrus"
	"net"
	"strings"
)

type tcpAcceptor struct {
	peer.SessionManager
	peer.CorePeerProperty
	peer.CoreContextSet
	peer.CoreRunningTag
	peer.CoreProcBundle
	peer.CoreTCPSocketOption
	peer.CoreCaptureIOPanic

	listener net.Listener
}

func (self *tcpAcceptor) Port() int {
	if self.listener == nil {
		return 0
	}

	return self.listener.Addr().(*net.TCPAddr).Port
}

func (self *tcpAcceptor) IsReady() bool {
	return self.IsRunning()
}

func (self *tcpAcceptor) Start() gorpc.Peer {
	self.WaitStopFinished()
	if self.IsRunning() {
		return self
	}
	ln, err := util.DetectPort(self.Address(), func(a *util.Address, port int) (interface{}, error) {
		return net.Listen("tcp", a.HostPortString(port))
	})
	if err != nil {
		logrus.Errorf("#tcp.listen failed(%s) %v", self.Name(), err.Error())
		self.SetRunning(false)
		return self
	}

	self.listener = ln.(net.Listener)
	logrus.Infof("#tcp.listen(%s) %s", self.Name(), self.ListenAddress())
	go self.accept()
	return self
}

func (self *tcpAcceptor) ListenAddress() string {
	pos := strings.Index(self.Address(), ":")
	if pos == -1 {
		return self.Address()
	}

	host := self.Address()[:pos]
	return util.JoinAddress(host, self.Port())
}

func (self *tcpAcceptor) accept() {
	self.SetRunning(true)

	for {
		conn, err := self.listener.Accept()
		if self.IsStoping() {
			break
		}

		if err != nil {
			logrus.Errorf("#tcp.accept failed(%s) %v", self.Name(), err.Error())
			break
		}
		go self.onNewSession(conn)
	}
	self.SetRunning(false)
	self.EndStopping()
}

func (self *tcpAcceptor) onNewSession(conn net.Conn) {
	self.ApplySocketOption(conn)
	ses := newSession(conn, self, nil)
	ses.Start()
	self.ProcEvent(&gorpc.RecvMsgEvent{
		Ses: ses,
		Msg: &gorpc.SessionAccepted{},
	})
}

func (a *tcpAcceptor) Stop() {
	if !a.IsRunning() {
		return
	}

	if a.IsStoping() {
		return
	}

	a.StartStoping()
	a.listener.Close()
	a.CloseAllSession()
	a.WaitStopFinished()
}

func (a *tcpAcceptor) TypeName() string {
	return "tcp.Acceptor"
}

func init() {
	peer.RegisterPeerCreator(func() gorpc.Peer {
		p := &tcpAcceptor{
			SessionManager: new(peer.CoreSessionManager),
		}
		p.CoreTCPSocketOption.Init()
		return p
	})
}
