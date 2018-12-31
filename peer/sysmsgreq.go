package peer

import (
	"github.com/0990/gorpc"
	"github.com/0990/gorpc/codec"
	_ "github.com/0990/gorpc/codec/binary"
	"github.com/0990/gorpc/util"
	"reflect"
)

func init() {
	gorpc.RegisterMessageMeta(&gorpc.MessageMeta{
		Codec: codec.MustGetCodec("binary"),
		Type:  reflect.TypeOf((*gorpc.SessionAccepted)(nil)).Elem(),
		ID:    int(util.StringHash("gorpc.SessionAccepted")),
	})

	gorpc.RegisterMessageMeta(&gorpc.MessageMeta{
		Codec: codec.MustGetCodec("binary"),
		Type:  reflect.TypeOf((*gorpc.SessionConnected)(nil)).Elem(),
		ID:    int(util.StringHash("gorpc.SessionConnected")),
	})
	gorpc.RegisterMessageMeta(&gorpc.MessageMeta{
		Codec: codec.MustGetCodec("binary"),
		Type:  reflect.TypeOf((*gorpc.SessionConnectError)(nil)).Elem(),
		ID:    int(util.StringHash("gorpc.SessionConnectError")),
	})
	gorpc.RegisterMessageMeta(&gorpc.MessageMeta{
		Codec: codec.MustGetCodec("binary"),
		Type:  reflect.TypeOf((*gorpc.SessionClosed)(nil)).Elem(),
		ID:    int(util.StringHash("gorpc.SessionClosed")),
	})
	gorpc.RegisterMessageMeta(&gorpc.MessageMeta{
		Codec: codec.MustGetCodec("binary"),
		Type:  reflect.TypeOf((*gorpc.SessionCloseNotify)(nil)).Elem(),
		ID:    int(util.StringHash("gorpc.SessionCloseNotify")),
	})
	gorpc.RegisterMessageMeta(&gorpc.MessageMeta{
		Codec: codec.MustGetCodec("binary"),
		Type:  reflect.TypeOf((*gorpc.SessionInit)(nil)).Elem(),
		ID:    int(util.StringHash("gorpc.SessionInit")),
	})
}
