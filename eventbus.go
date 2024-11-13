package eventbus

import (
	"context"
	"sync"
)

type (
	Publisher interface {
		Publish(ctx context.Context, event Event)
		Flush(ctx context.Context, events *Events)
		Wait()
	}
	Subscriber interface {
		Subscribe(name EventName, handler Handler, options ...HandlerOption) (HandlerID, func())
		Unsubscribe(name EventName, id HandlerID)
	}
)

type EventBus struct {
	handlers map[EventName]handlerOptions
	nextID   HandlerID
	mu       sync.RWMutex
	wg       sync.WaitGroup
}

func New() *EventBus {
	return &EventBus{
		handlers: make(map[EventName]handlerOptions),
		nextID:   1,
	}
}

func (e *EventBus) Subscribe(name EventName, handler Handler, options ...HandlerOption) (HandlerID, func()) {
	e.mu.Lock()
	defer e.mu.Unlock()

	handlerID := e.nextID
	e.nextID++

	if e.handlers[name] == nil {
		e.handlers[name] = make(handlerOptions)
	}
	e.handlers[name][handlerID] = newHandlerOption(handler, &e.wg, options)

	return handlerID, func() {
		e.Unsubscribe(name, handlerID)
	}
}

func (e *EventBus) Unsubscribe(name EventName, id HandlerID) {
	e.mu.Lock()
	defer e.mu.Unlock()

	handlers, ok := e.handlers[name]
	if !ok {
		return
	}

	delete(handlers, id)

	if len(handlers) == 0 {
		delete(e.handlers, name)
	}
}

func (e *EventBus) Flush(ctx context.Context, events *Events) {
	for _, event := range events.Release() {
		e.Publish(ctx, event)
	}
}

func (e *EventBus) Publish(ctx context.Context, event Event) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	handlers, ok := e.handlers[event.Name()]
	if !ok {
		return
	}

	for _, handler := range handlers {
		handler.Handle(ctx, event)
	}
}

func (e *EventBus) Wait() {
	e.wg.Wait()
}
