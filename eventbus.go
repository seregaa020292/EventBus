package eventbus

import (
	"context"
	"sync"
)

type (
	Publisher interface {
		Publish(ctx context.Context, event Event)
		Flush(ctx context.Context, events *Events)
	}
	Subscriber interface {
		Subscribe(name EventName, handler Handler, options ...HandlerOption) (HandlerID, func())
		Unsubscribe(name EventName, id HandlerID)
	}
)

type EventBus struct {
	handlers map[EventName]Handlers
	nextID   HandlerID
	mu       sync.RWMutex
}

func New() *EventBus {
	return &EventBus{
		handlers: make(map[EventName]Handlers),
		nextID:   1,
	}
}

func (e *EventBus) Subscribe(name EventName, handler Handler, options ...HandlerOption) (HandlerID, func()) {
	e.mu.Lock()
	defer e.mu.Unlock()

	handlerWithOption := newHandlerOption(handler)
	for _, option := range options {
		option(handlerWithOption)
	}

	handlerID := e.nextID
	e.nextID++

	if e.handlers[name] == nil {
		e.handlers[name] = make(Handlers)
	}
	e.handlers[name][handlerID] = handlerWithOption

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
		h := handler.(*handlerOption)
		if h.isAsync {
			go h.Handle(ctx, event)
		} else {
			h.Handle(ctx, event)
		}
	}
}
