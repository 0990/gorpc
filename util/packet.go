package util

import (
	"encoding/binary"
	"errors"
	"github.com/0990/gorpc"
	"github.com/0990/gorpc/codec"
	"io"
)

var (
	ErrMaxPacket  = errors.New("packet over size")
	ErrMinPacket  = errors.New("packet short size")
	ErrShortMsgID = errors.New("short msgid")
)

const (
	bodySize  = 2
	msgIDSize = 2
)

func RecvLTVPacket(reader io.Reader, maxPacketSize int) (msg interface{}, err error) {
	var sizeBuffer = make([]byte, bodySize)

	_, err = io.ReadFull(reader, sizeBuffer)

	if err != nil {
		return
	}

	if len(sizeBuffer) < bodySize {
		return nil, ErrMinPacket
	}

	size := binary.LittleEndian.Uint16(sizeBuffer)

	if maxPacketSize > 0 && size >= uint16(maxPacketSize) {
		return nil, ErrMaxPacket
	}

	body := make([]byte, size)
	_, err = io.ReadFull(reader, body)
	if err != nil {
		return
	}

	if len(body) < bodySize {
		return nil, ErrShortMsgID
	}

	msgid := binary.LittleEndian.Uint16(body)
	msgData := body[msgIDSize:]
	msg, _, err = codec.DecodeMessage(int(msgid), msgData)
	if err != nil {
		return nil, err
	}
	return
}

func SendLTVPacket(writer io.Writer, data interface{}) error {
	var (
		msgData []byte
		msgID   int
		meta    *gorpc.MessageMeta
	)

	switch m := data.(type) {
	case *gorpc.RawPacket:
		msgData = m.MsgData
		msgID = m.MsgID
	default:
		var err error
		msgData, meta, err = codec.EncodeMessage(data)
		if err != nil {
			return err
		}
		msgID = meta.ID
	}

	pkt := make([]byte, bodySize+msgIDSize+len(msgData))

	binary.LittleEndian.PutUint16(pkt, uint16(msgIDSize+len(msgData)))
	binary.LittleEndian.PutUint16(pkt[msgIDSize:], uint16(msgID))

	copy(pkt[bodySize+msgIDSize:], msgData)

	err := WriteFull(writer, pkt)

	if meta != nil {
		codec.FreeCodecResource(meta.Codec, msgData)
	}
	return err
}
