package gorpc

type Peer interface {
	Start() Peer
	Stop()
	TypeName() string
}

type PeerProperty interface {
	Name() string
	Address() string
	SetName(v string)
	SetAddress(v string)
	SetQueue(v EventQueue)
}

type GenericPeer interface {
	Peer
	PeerProperty
}

type ContextSet interface {
	SetContext(key interface{}, v interface{})
	GetContext(key interface{}) (interface{}, bool)
	FetchContext(key, valuePtr interface{}) bool
}

type SessionAccessor interface {
	GetSession(int64) Session

	VisitSession(func(Session) bool)
	SessionCount() int
	CloseAllSession()
}

type PeerReadyChecker interface {
	IsReady() bool
}

type PeerCaptureIOPanic interface {
	EnableCaptureIOPanic(v bool)
	CaptureIOPanic() bool
}
