package tcp

import (
	"github.com/0990/gorpc"
	"github.com/0990/gorpc/peer"
	"net"
	"sync"
	"time"
)

type tcpConnector struct {
	peer.SessionManager
	peer.CorePeerProperty
	peer.CoreContextSet
	peer.CoreRunningTag
	peer.CoreProcBundle
	peer.CoreTCPSocketOption

	defautSes *tcpSession

	tryConnTimes int
	sesEndSignal sync.WaitGroup

	reconDur time.Duration
}

func (c *tcpConnector) Start() gorpc.Peer {
	c.WaitStopFinished()

	if c.IsRunning() {
		return c
	}

	go c.connect(c.Address())

	return c
}

func (c *tcpConnector) Session() gorpc.Session {
	return c.defautSes
}

func (c *tcpConnector) SetSessionManager(raw interface{}) {
	c.SessionManager = raw.(peer.SessionManager)
}

func (c *tcpConnector) Stop() {
	if !c.IsRunning() {
		return
	}

	if c.IsStoping() {
		return
	}

	c.StartStoping()

	c.defautSes.Close()
	c.WaitStopFinished()
}

func (c *tcpConnector) ReconnectDuration() time.Duration {
	return c.reconDur
}

func (c *tcpConnector) SetReconnectDuration(v time.Duration) {
	c.reconDur = v
}

func (c *tcpConnector) Port() int {
	if c.defautSes.conn == nil {
		return 0
	}

	return c.defautSes.conn.LocalAddr().(*net.TCPAddr).Port
}

const reportConnectFailedLimitTimes = 3

func (c *tcpConnector) connect(address string) {
	c.SetRunning(true)
	for {
		c.tryConnTimes++
		conn, err := net.Dial("tcp", address)
		c.defautSes.conn = conn
		if err != nil {
			if c.tryConnTimes <= reportConnectFailedLimitTimes {
				log.Errorf("#tcp.connect failed(%s) %v", c.Name(), err.Error())
				if c.tryConnTimes == reportConnectFailedLimitTimes {
					log.Errorf("(%s) continue reconnecting,but mute log", c.Name())
				}
			}

			if c.ReconnectDuration() == 0 || c.IsStoping() {
				c.ProcEvent(&gorpc.RecvMsgEvent{
					Ses: c.defautSes,
					Msg: &gorpc.SessionConnectError{},
				})
				break
			}
			time.Sleep(c.ReconnectDuration())
			continue
		}
		c.sesEndSignal.Add(1)
		c.ApplySocketOption(conn)
		c.defautSes.Start()
		c.tryConnTimes = 0

		c.ProcEvent(&gorpc.RecvMsgEvent{
			Ses: c.defautSes,
			Msg: &gorpc.SessionConnected{},
		})

		c.sesEndSignal.Wait()
		c.defautSes.conn = nil
		if c.IsStoping() || c.ReconnectDuration() == 0 {
			break
		}

		time.Sleep(c.ReconnectDuration())
		continue
	}

	c.SetRunning(false)
	c.EndStopping()
}

func (c *tcpConnector) IsReady() bool {
	return c.SessionCount() != 0
}

func (c *tcpConnector) TypeName() string {
	return "tcp.Connector"
}

func init() {
	peer.RegisterPeerCreator(func() gorpc.Peer {
		c := &tcpConnector{
			SessionManager: new(peer.CoreSessionManager),
		}

		c.defautSes = newSession(nil, c, func() {
			c.sesEndSignal.Done()
		})
		c.CoreTCPSocketOption.Init()
		return c
	})
}
