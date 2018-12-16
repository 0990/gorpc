package rpc

import (
	"github.com/0990/gorpc"
	"github.com/0990/gorpc/codec"
	_ "github.com/0990/gorpc/codec/protoplus"
	"github.com/davyxu/protoplus/proto"
	"reflect"
)

type RemoteCallREQ struct {
	MsgID  uint32
	Data   []byte
	CallID int64
}

func (self RemoteCallREQ) String() string {
	return proto.CompactTextString(self)
}

func (self *RemoteCallREQ) Size() (ret int) {
	ret += proto.SizeUInt32(0, self.MsgID)
	ret += proto.SizeBytes(1, self.Data)
	ret += proto.SizeInt64(2, self.CallID)
	return
}

func (self *RemoteCallREQ) Marshal(buffer *proto.Buffer) error {
	proto.MarshalUInt32(buffer, 0, self.MsgID)
	proto.MarshalBytes(buffer, 1, self.Data)
	proto.MarshalInt64(buffer, 2, self.CallID)
	return nil
}

func (self *RemoteCallREQ) Unmarshal(buffer *proto.Buffer, fieldIndex uint64, wt proto.WireType) error {
	switch fieldIndex {
	case 0:
		return proto.UnmarshalUInt32(buffer, wt, &self.MsgID)
	case 1:
		return proto.UnmarshalBytes(buffer, wt, &self.Data)
	case 2:
		return proto.UnmarshalInt64(buffer, wt, &self.CallID)
	}
	return proto.ErrUnknownField
}

type RemoteCallACK struct {
	MsgID  uint32
	Data   []byte
	CallID int64
}

func (self *RemoteCallACK) String() string {
	return proto.CompactTextString(self)
}

func (self *RemoteCallACK) Marshal(buffer *proto.Buffer) error {
	proto.MarshalUInt32(buffer, 0, self.MsgID)
	proto.MarshalBytes(buffer, 1, self.Data)
	proto.MarshalInt64(buffer, 2, self.CallID)
	return nil
}

func (self *RemoteCallACK) Unmarshal(buffer *proto.Buffer, fieldIndex uint64, wt proto.WireType) error {
	switch fieldIndex {
	case 0:
		return proto.UnmarshalUInt32(buffer, wt, &self.MsgID)
	case 1:
		return proto.UnmarshalBytes(buffer, wt, &self.Data)
	case 2:
		return proto.UnmarshalInt64(buffer, wt, &self.CallID)
	}
	return proto.ErrUnknownField
}

func init() {
	gorpc.RegisterMessageMeta(&gorpc.MessageMeta{
		Codec: codec.MustGetCodec("protoplus"),
		Type:  reflect.TypeOf((*RemoteCallREQ)(nil)).Elem(),
		ID:    58654,
	})
	gorpc.RegisterMessageMeta(&gorpc.MessageMeta{
		Codec: codec.MustGetCodec("protoplus"),
		Type:  reflect.TypeOf((*RemoteCallACK)(nil)).Elem(),
		ID:    20476,
	})
}
