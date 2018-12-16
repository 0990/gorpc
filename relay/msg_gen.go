package relay

import (
	"github.com/0990/gorpc"
	"github.com/0990/gorpc/codec"
	_ "github.com/0990/gorpc/codec/protoplus"
	"github.com/davyxu/protoplus/proto"
	"reflect"
)

type RelayACK struct {
	Msg        []byte
	MsgID      uint32
	Bytes      []byte
	Int64      int64
	Int64Slice []int64
	Str        string
}

func (self *RelayACK) String() string {
	return proto.CompactTextString(self)
}

func (self *RelayACK) Size() (ret int) {
	ret += proto.SizeBytes(0, self.Msg)
	ret += proto.SizeUInt32(1, self.MsgID)
	ret += proto.SizeBytes(2, self.Bytes)
	ret += proto.SizeInt64(3, self.Int64)
	ret += proto.SizeInt64Slice(4, self.Int64Slice)
	ret += proto.SizeString(5, self.Str)
	return
}

func (self *RelayACK) Marshal(buffer *proto.Buffer) error {
	proto.MarshalBytes(buffer, 0, self.Msg)
	proto.MarshalUInt32(buffer, 1, self.MsgID)
	proto.MarshalBytes(buffer, 2, self.Bytes)
	proto.MarshalInt64(buffer, 3, self.Int64)
	proto.MarshalInt64Slice(buffer, 4, self.Int64Slice)
	proto.MarshalString(buffer, 5, self.Str)
	return nil
}

func (self *RelayACK) Unmarshal(buffer *proto.Buffer, fieldIndex uint64, wt proto.WireType) error {
	switch fieldIndex {
	case 0:
		return proto.UnmarshalBytes(buffer, wt, &self.Msg)
	case 1:
		return proto.UnmarshalUInt32(buffer, wt, &self.MsgID)
	case 2:
		return proto.UnmarshalBytes(buffer, wt, &self.Bytes)
	case 3:
		return proto.UnmarshalInt64(buffer, wt, &self.Int64)
	case 4:
		return proto.UnmarshalInt64Slice(buffer, wt, &self.Int64Slice)
	case 5:
		return proto.UnmarshalString(buffer, wt, &self.Str)

	}

	return proto.ErrUnknownField
}

func init() {
	gorpc.RegisterMessageMeta(&gorpc.MessageMeta{
		Codec: codec.MustGetCodec("protoplus"),
		Type:  reflect.TypeOf((*RelayACK)(nil)).Elem(),
		ID:    45545,
	})
}
