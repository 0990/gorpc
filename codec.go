package gorpc

type Codec interface {
	Encode(msgObj interface{}) (data interface{}, err error)
	Decode(data interface{}, msgObj interface{}) error
	Name() string
	MimeType() string
}
