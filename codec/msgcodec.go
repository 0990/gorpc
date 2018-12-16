package codec

import "github.com/0990/gorpc"

func EncodeMessage(msg interface{}) (data []byte, meta *gorpc.MessageMeta, err error) {
	meta = gorpc.MessageMetaByMsg(msg)
	if meta == nil {
		return nil, nil, gorpc.NewErrContext("msg not exists", msg)
	}

	var raw interface{}
	raw, err = meta.Codec.Encode(msg)
	if err != nil {
		return
	}
	data = raw.([]byte)
	return
}

func DecodeMessage(msgid int, data []byte) (interface{}, *gorpc.MessageMeta, error) {
	meta := gorpc.MessageMetaByID(msgid)

	if meta == nil {
		return nil, nil, gorpc.NewErrContext("msg not exists", msgid)
	}

	msg := meta.NewType()

	err := meta.Codec.Decode(data, msg)
	if err != nil {
		return nil, meta, err
	}
	return msg, meta, nil
}

type CodeRecycler interface {
	Free(data interface{})
}

func FreeCodecResource(codec gorpc.Codec, data interface{}) {
	if codec == nil {
		return
	}

	if recycler, ok := codec.(CodeRecycler); ok {
		recycler.Free(data)
	}
}
