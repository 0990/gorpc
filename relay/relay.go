package relay

import (
	"errors"
	"github.com/0990/gorpc"
	"github.com/0990/gorpc/codec"
	"github.com/sirupsen/logrus"
)

var (
	ErrInvalidPeerSession = errors.New("Require valid gorpc.Session or gorpc.TCPConnector")
)

func Relay(sesDetector interface{}, dataList ...interface{}) error {
	ses, err := getSession(sesDetector)
	if err != nil {
		logrus.Errorln("relay.Relay:", err)
		return err
	}

	var ack RelayACK
	for _, payload := range dataList {
		switch value := payload.(type) {
		case int64:
			ack.Int64 = value
		case []int64:
			ack.Int64Slice = value

		case string:
			ack.Str = value
		case []byte: // 作为payload
			ack.Bytes = value
		default:
			if ack.MsgID == 0 {
				var meta *gorpc.MessageMeta
				ack.Msg, meta, err = codec.EncodeMessage(payload)

				if err != nil {
					return err
				}

				ack.MsgID = uint32(meta.ID)
			} else {
				panic("Multi message relay not support")
			}

		}
	}
	ses.Send(&ack)
	return nil
}

func getSession(sesDetector interface{}) (gorpc.Session, error) {
	switch unknown := sesDetector.(type) {
	case gorpc.Session:
		return unknown, nil
	case gorpc.TCPConnector:
		return unknown.Session(), nil
	default:
		return nil, ErrInvalidPeerSession
	}
}
