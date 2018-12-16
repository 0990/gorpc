package codec

import "github.com/0990/gorpc"

var registedCodes []gorpc.Codec

func RegisterCodec(c gorpc.Codec) {
	if GetCodec(c.Name()) != nil {
		panic("duplicate codec:" + c.Name())
	}

	registedCodes = append(registedCodes, c)
}

func GetCodec(name string) gorpc.Codec {
	for _, c := range registedCodes {
		if c.Name() == name {
			return c
		}
	}
	return nil
}

func MustGetCodec(name string) gorpc.Codec {
	codec := GetCodec(name)
	if codec == nil {
		panic("codec not register!" + name)
	}
	return codec
}
