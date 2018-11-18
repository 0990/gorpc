package gorpc

import (
	"log"
	"runtime/debug"
	"sync"
)

type EventQueue interface {
	StartLoop() EventQueue
	StopLoop() EventQueue
	Wait()
	Post(callback func())
	EnableCapturePanic(v bool)
}

type eventQueue struct {
	*Pipe
	endSignal    sync.WaitGroup
	capturePanic bool
}

func (self *eventQueue) EnableCapturePanic(v bool) {
	self.capturePanic = v
}

func (self *eventQueue) Post(callback func()) {
	if callback == nil {
		return
	}
	self.Add(callback)
}

func (self *eventQueue) protectedCall(callback func()) {
	if self.capturePanic {
		defer func() {
			if err := recover(); err != nil {
				debug.PrintStack()
			}
		}()
	}
	callback()
}

func (self *eventQueue) StartLoop() EventQueue {
	self.endSignal.Add(1)
	go func() {
		var writeList []interface{}

		for {
			writeList = writeList[0:0]
			exit := self.Pick(&writeList)

			for _, msg := range writeList {
				switch t := msg.(type) {
				case func():
					self.protectedCall(t)
				case nil:
					break
				default:
					log.Printf("unexpected type:%T", t)
				}
			}
			if exit {
				break
			}
		}
		self.endSignal.Done()
	}()
	return self
}

func (self *eventQueue) StopLoop() EventQueue {
	self.Add(nil)
	return self
}

func (self *eventQueue) Wait() {
	self.endSignal.Wait()
}

func NewEventQueue() EventQueue {
	return &eventQueue{
		Pipe: NewPipe(),
	}
}

func SessionQueuedCall(ses Session, callback func()) {
	if ses == nil {
		return
	}

	q := ses.Peer().(interface {
		Queue() EventQueue
	}).Queue()
	QueuedCall(q, callback)
}

func QueuedCall(queue EventQueue, callback func()) {
	if queue == nil {
		callback()
	} else {
		queue.Post(callback)
	}
}
