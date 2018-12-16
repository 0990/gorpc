package tests

import (
	"fmt"
	"github.com/0990/gorpc"
	"github.com/0990/gorpc/codec"
	_ "github.com/0990/gorpc/codec/binary"
	"github.com/0990/gorpc/util"
	"reflect"
)

type TestEchoACK struct {
	Msg   string
	Value int32
}

func (self *TestEchoACK) String() string { return fmt.Sprintf("%+v", *self) }

func init() {
	gorpc.RegisterMessageMeta(&gorpc.MessageMeta{
		Codec: codec.MustGetCodec("binary"),
		Type:  reflect.TypeOf((*TestEchoACK)(nil)).Elem(),
		ID:    int(util.StringHash("tests.TestEchoACK")),
	})
}
