package tcp

import (
	"github.com/0990/gorpc"
	"github.com/0990/gorpc/util"
	"io"
	"net"
)

type TCPMessageTransmitter struct {
}

type socketOpt interface {
	MaxPacketSize() int
	ApplySocketReadTimeout(conn net.Conn, callback func())
	ApplySocketWriteTimeout(conn net.Conn, callback func())
}

func (TCPMessageTransmitter) OnRecvMessage(ses gorpc.Session) (msg interface{}, err error) {
	reader, ok := ses.Raw().(io.Reader)

	if !ok || reader == nil {
		return nil, nil
	}

	opt, ok := ses.Peer().(socketOpt)
	if !ok {
		panic("hello")
	}

	if conn, ok := ses.Raw().(net.Conn); ok {
		opt.ApplySocketReadTimeout(conn, func() {
			msg, err = util.RecvLTVPacket(reader, opt.MaxPacketSize())
		})
	}
	return
}

func (TCPMessageTransmitter) OnSendMessage(ses gorpc.Session, msg interface{}) (err error) {
	writer, ok := ses.Raw().(io.Writer)
	if !ok || writer == nil {
		return nil
	}

	opt, ok := ses.Peer().(socketOpt)
	if !ok {
		panic("hello")
	}
	opt.ApplySocketWriteTimeout(ses.Raw().(net.Conn), func() {
		err = util.SendLTVPacket(writer, msg)
	})
	return
}
