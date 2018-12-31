package proc

import (
	"github.com/0990/gorpc"
	"reflect"
	"sync"
)

type SyncReceiver struct {
	evChan   chan gorpc.Event
	callback func(ev gorpc.Event)
}

func (p *SyncReceiver) EventCallback() gorpc.EventCallback {
	return p.callback
}

func (p *SyncReceiver) Recv(callback gorpc.EventCallback) *SyncReceiver {
	callback(<-p.evChan)
	return p
}

func (p *SyncReceiver) WaitMessage(msgName string) (msg interface{}) {
	var wg sync.WaitGroup

	meta := gorpc.MessageMetaByFullName(msgName)
	if meta == nil {
		panic("unknow message name:" + msgName)
	}
	wg.Add(1)

	p.Recv(func(ev gorpc.Event) {
		inMeta := gorpc.MessageMetaByType(reflect.TypeOf(ev.Message()))
		if inMeta == meta {
			msg = ev.Message()
			wg.Done()
		}
	})
	wg.Wait()
	return
}

func NewSyncReceiver() *SyncReceiver {
	sr := &SyncReceiver{
		evChan: make(chan gorpc.Event),
	}
	sr.callback = func(ev gorpc.Event) {
		sr.evChan <- ev
	}

	return sr
}
