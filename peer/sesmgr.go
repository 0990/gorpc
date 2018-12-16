package peer

import (
	"github.com/0990/gorpc"
	"sync"
	"sync/atomic"
)

type SessionManager interface {
	gorpc.SessionAccessor
	Add(session gorpc.Session)
	Remove(session gorpc.Session)
	Count() int
	SetIDBase(base int64)
}

type CoreSessionManager struct {
	sesById  sync.Map
	sesIDGen int64
	count    int64
}

func (self *CoreSessionManager) SetIDBase(base int64) {
	atomic.StoreInt64(&self.sesIDGen, base)
}

func (self *CoreSessionManager) Count() int {
	return int(atomic.LoadInt64(&self.count))
}

func (self *CoreSessionManager) Add(ses gorpc.Session) {
	id := atomic.AddInt64(&self.sesIDGen, 1)
	atomic.AddInt64(&self.count, 1)
	ses.(interface {
		SetID(int64)
	}).SetID(id)

	self.sesById.Store(id, ses)
}

func (self *CoreSessionManager) Remove(ses gorpc.Session) {
	self.sesById.Delete(ses.ID())
	atomic.AddInt64(&self.count, -1)
}

func (self *CoreSessionManager) GetSession(id int64) gorpc.Session {
	if v, ok := self.sesById.Load(id); ok {
		return v.(gorpc.Session)
	}
	return nil
}

func (self *CoreSessionManager) VisitSession(callback func(session gorpc.Session) bool) {
	self.sesById.Range(func(key, value interface{}) bool {
		return callback(value.(gorpc.Session))
	})
}

func (self *CoreSessionManager) CloseAllSession() {
	self.VisitSession(func(ses gorpc.Session) bool {
		ses.Close()
		return true
	})
}

func (self *CoreSessionManager) SessionCount() int {
	v := atomic.LoadInt64(&self.count)
	return int(v)
}
