package gorpc

type Event interface {
	Session() Session
	Message() interface{}
}

type MessageTransmitter interface {
	OnRecvMessage(ses Session) (msg interface{}, err error)
	OnSendMessage(ses Session, msg interface{}) error
}

type EventHooker interface {
	OnInboundEvent(input Event) (output Event)
	OnOutboundEvent(input Event) (output Event)
}

type EventCallback func(ev Event)
