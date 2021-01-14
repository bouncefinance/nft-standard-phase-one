package eventLoop

import (
	"context"
	logs "log"
	"reflect"
	"sync"
	"time"
)

const (
	cmdEmitEvent        = 0
	cmdContinue         = 1
	cmdClose            = 2
	cmdCloseImmediately = 3
)

type EventCallback func(*EventData) bool

type EventData struct {
	Event     string
	Timestamp int64
	Data      interface{}
}

type eventCommand struct {
	sort int
	data *EventData
}

type eventListener struct {
	sync.Mutex
	callback EventCallback
	isOnce   bool
	fired    bool
}

type EventLoop struct {
	sync.Mutex
	cmdCh    chan *eventCommand
	channels map[string][]*eventListener
	ctx      context.Context
	runLock  sync.Mutex
	running  int64
	limit    int64
	queue    []*eventCommand
}

func (el *EventLoop) On(eventName string, callback EventCallback) {
	el.add(eventName, callback, false)
}
func (el *EventLoop) Once(eventName string, callback EventCallback) {
	el.add(eventName, callback, true)
}
func (el *EventLoop) add(eventName string, callback EventCallback, isOnce bool) {
	el.Lock()
	defer el.Unlock()

	if callback == nil || len(eventName) == 0 {
		return
	}

	listener := &eventListener{
		callback: callback,
		isOnce:   isOnce,
		fired:    false,
	}

	channel := el.channels[eventName]
	if channel == nil {
		channel = make([]*eventListener, 1)
		channel[0] = listener
	} else {
		channel = append(channel, listener)
	}
	el.channels[eventName] = channel
}
func (el *EventLoop) Off(eventName string, callback EventCallback) {
	el.Lock()
	defer el.Unlock()

	channel := el.channels[eventName]
	if channel == nil {
		return
	}

	ref := reflect.ValueOf(callback).Pointer()
	count := 0
	listeners := make([]*eventListener, len(channel))
	for _, listener := range channel {
		if listener.fired {
			continue
		}
		cref := reflect.ValueOf(listener.callback).Pointer()
		if cref != ref {
			listeners[count] = listener
			count++
		}
	}
	ls := make([]*eventListener, count)
	for i := 0; i < count; i++ {
		ls[i] = listeners[i]
	}
	el.channels[eventName] = ls
	listeners = listeners[:0]
}
func (el *EventLoop) OffAll(eventName string) {
	el.Lock()
	defer el.Unlock()

	channel := el.channels[eventName]
	if channel == nil {
		return
	}

	el.channels[eventName] = make([]*eventListener, 0)
}
func (el *EventLoop) Count(eventName string) uint64 {
	el.Lock()
	defer el.Unlock()

	channel := el.channels[eventName]
	if channel == nil {
		return uint64(0)
	} else {
		return uint64(len(channel))
	}
}
func (el *EventLoop) Emit(eventName string, data interface{}) {
	edata := &EventData{
		Event:     eventName,
		Timestamp: time.Now().UnixNano(),
		Data:      data,
	}
	command := &eventCommand{
		sort: cmdEmitEvent,
		data: edata,
	}
	el.cmdCh <- command
}
func (el *EventLoop) EmitSync(eventName string, data interface{}) {
	if len(eventName) == 0 {
		return
	}

	edata := &EventData{
		Event:     eventName,
		Timestamp: time.Now().UnixNano(),
		Data:      data,
	}
	cmd := &eventCommand{
		sort: cmdEmitEvent,
		data: edata,
	}

	el.Lock()
	channel := el.channels[eventName]
	if channel == nil {
		el.Unlock()
		return
	}

	copy := make([]*eventListener, len(channel))
	i := 0
	for _, l := range channel {
		l.Lock()
		if !l.fired {
			copy[i] = l
			i++
		}
		l.Unlock()
	}
	channel = copy[:i]
	el.channels[eventName] = channel

	copy = make([]*eventListener, len(channel))
	for i, c := range channel {
		copy[i] = c
	}
	el.Unlock()

	for _, l := range copy {
		l.Lock()
		cb := l.callback
		if l.fired {
			l.Unlock()
			continue
		}
		l.fired = l.isOnce
		l.Unlock()
		if cb(cmd.data) {
			break
		}
	}
}
func (el *EventLoop) Close() {
	command := &eventCommand{
		sort: cmdClose,
	}
	el.cmdCh <- command
}
func (el *EventLoop) CloseImmediately() {
	command := &eventCommand{
		sort: cmdCloseImmediately,
	}
	el.cmdCh <- command
}
func (el *EventLoop) HasEvent(eventName string) bool {
	el.Lock()
	defer el.Unlock()
	_, ok := el.channels[eventName]
	return ok
}
func (el *EventLoop) start() {
loop:
	for {
		select {
		case <-el.ctx.Done():
			return
		case cmd := <-el.cmdCh:
			switch cmd.sort {
			case cmdClose:
				break loop
			case cmdCloseImmediately:
				break loop
			case cmdEmitEvent:
				if cmd.data == nil || len(cmd.data.Event) == 0 {
					continue loop
				}
				el.runLock.Lock()
				if el.limit > 0 && el.running >= el.limit {
					el.queue = append(el.queue, cmd)
					el.runLock.Unlock()
					continue loop
				}
				el.runLock.Unlock()

				el.Lock()
				channel := el.channels[cmd.data.Event]
				if channel == nil {
					el.Unlock()
					continue loop
				}

				copy := make([]*eventListener, len(channel))
				i := 0
				for _, l := range channel {
					l.Lock()
					if !l.fired {
						copy[i] = l
						i++
					}
					l.Unlock()
				}
				channel = copy[:i]
				el.channels[cmd.data.Event] = channel
				copy = make([]*eventListener, len(channel))
				for i, c := range channel {
					copy[i] = c
				}
				el.Unlock()

				go func(channel []*eventListener) {
					el.runLock.Lock()
					el.running++
					el.runLock.Unlock()
					for _, l := range channel {
						l.Lock()
						cb := l.callback
						if l.fired {
							l.Unlock()
							continue
						}
						l.fired = l.isOnce
						l.Unlock()
						if cb(cmd.data) {
							break
						}
					}
					el.runLock.Lock()
					el.running--
					el.runLock.Unlock()
					command := &eventCommand{
						sort: cmdContinue,
					}
					el.cmdCh <- command
				}(copy)
			case cmdContinue:
				el.runLock.Lock()
				if (el.limit > 0 && el.running >= el.limit) || (len(el.queue) == 0) {
					el.runLock.Unlock()
					continue loop
				}
				cmd = el.queue[0]
				el.queue = el.queue[1:]
				el.runLock.Unlock()
				go func() {
					el.cmdCh <- cmd
				}()
			}
		}
	}
}

func New(ctx context.Context) *EventLoop {
	el := &EventLoop{
		cmdCh:    make(chan *eventCommand),
		channels: make(map[string][]*eventListener),
		ctx:      ctx,
	}
	go el.start()
	return el
}
func NewWithLimit(ctx context.Context, limit int64) *EventLoop {
	el := &EventLoop{
		cmdCh:    make(chan *eventCommand),
		channels: make(map[string][]*eventListener),
		ctx:      ctx,
		limit:    limit,
	}
	go el.start()
	return el
}
