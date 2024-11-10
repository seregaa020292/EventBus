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
		Subscribe(typ EventType, handler Handler) (HandlerID, func())
		Unsubscribe(typ EventType, id HandlerID)
	}
)

type EventBus struct {
	handlers map[EventType]Handlers
	nextID   uint16
	mu       sync.RWMutex
}

func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[EventType]Handlers),
		nextID:   1,
	}
}

func (e *EventBus) Subscribe(typ EventType, handler Handler) (HandlerID, func()) {
	e.mu.Lock()
	defer e.mu.Unlock()

	handlerID := e.nextID
	e.nextID++

	if e.handlers[typ] == nil {
		e.handlers[typ] = make(Handlers)
	}
	e.handlers[typ][handlerID] = handler

	unsubscribe := func() {
		e.Unsubscribe(typ, handlerID)
	}

	return handlerID, unsubscribe
}

func (e *EventBus) Unsubscribe(typ EventType, id HandlerID) {
	e.mu.Lock()
	defer e.mu.Unlock()

	handlers, ok := e.handlers[typ]
	if !ok {
		return
	}

	delete(handlers, id)

	if len(handlers) == 0 {
		delete(e.handlers, typ)
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

	handlers, ok := e.handlers[event.typ]
	if !ok {
		return
	}

	if event.isAsync {
		go e.publish(ctx, handlers, event.payload)
		return
	}

	e.publish(ctx, handlers, event.payload)
}

func (e *EventBus) publish(ctx context.Context, handlers Handlers, payload any) {
	for _, handler := range handlers {
		handler.Handle(ctx, payload)
	}
}
