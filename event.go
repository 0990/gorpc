package gorpc

type RecvMsgEvent struct {
	Ses Session
	Msg interface{}
}

func (rm *RecvMsgEvent) Session() Session {
	return rm.Ses
}

func (rm *RecvMsgEvent) Message() interface{} {
	return rm.Msg
}

func (rm *RecvMsgEvent) Send(msg interface{}) {
	rm.Ses.Send(msg)
}

func (rm *RecvMsgEvent) Reply(msg interface{}) {
	rm.Ses.Send(msg)
}

type SendMsgEvent struct {
	Ses Session
	Msg interface{}
}

func (sm *SendMsgEvent) Message() interface{} {
	return sm.Msg
}

func (sm *SendMsgEvent) Session() Session {
	return sm.Ses
}
