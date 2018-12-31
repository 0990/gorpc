package rpc

func (p *RemoteCallREQ) GetMsgID() uint16 {
	return uint16(p.MsgID)
}

func (p *RemoteCallREQ) GetMsgData() []byte {
	return p.Data
}

func (p *RemoteCallREQ) GetCallID() int64 {
	return p.CallID
}

func (p *RemoteCallACK) GetMsgID() uint16 {
	return uint16(p.MsgID)
}

func (p *RemoteCallACK) GetMsgData() []byte {
	return p.Data
}

func (p *RemoteCallACK) GetCallID() int64 {
	return p.CallID
}
