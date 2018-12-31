package gorpc

type SessionInit struct {
}

type SessionAccepted struct {
}

type SessionConnected struct {
}

type SessionConnectError struct {
}
type CloseReason int32

const (
	CloseReason_IO     CloseReason = iota // 普通IO断开
	CloseReason_Manual                    // 关闭前，调用过Session.Close
)

func (self CloseReason) String() string {
	switch self {
	case CloseReason_IO:
		return "IO"
	case CloseReason_Manual:
		return "Manual"
	}

	return "Unknown"
}

type SessionClosed struct {
	Reason CloseReason
}

type SessionCloseNotify struct {
}
