package protoplus

import (
	"github.com/0990/gorpc/codec"
	"github.com/davyxu/protoplus/proto"
)

type protoplus struct {
}

func (self *protoplus) Name() string {
	return "protoplus"
}

func (self *protoplus) MimeType() string {
	return "application/binary"
}

func (self *protoplus) Encode(msgObj interface{}) (data interface{}, err error) {
	return proto.Marshal(msgObj)
}

func (self *protoplus) Decode(data interface{}, msgObj interface{}) error {
	return proto.Unmarshal(data.([]byte), msgObj)
}

func init() {
	codec.RegisterCodec(new(protoplus))
}
