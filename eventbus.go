package eventbus

import (
	"context"
	"rand"
	"sync"
	"fmt"
)

type Publisher interface {
	Publish(ctx context.Context, event Event)
	Flush(ctx context.Context, queue *EventQueue)
	Wait()
}

type Subscriber interface {
	Subscribe(topic string, handler Handler, options ...HandlerOption) (string, func())
	Unsubscribe(topic string, id string)
}

type EventBus struct {
	config     Config
	handlers   map[string]map[string]*handler
	middleware chainMiddlewares
	mu         sync.RWMutex
	wg         sync.WaitGroup
}

func New(opts ...Option) *EventBus {
	return &EventBus{
		config:     newConfig(opts...),
		handlers:   make(map[string]map[string]*handler),
		middleware: middleware{middlewares: make([]Middleware, 0)},
	}
}

func (e *EventBus) Subscribe(topic string, h Handler, opts ...HandlerOption) (string, func()) {
	id := generateID()

	sub := &handler{
		base: h,
	}
	for _, opt := range opts {
		opt(sub)
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	if e.handlers[topic] == nil {
		e.handlers[topic] = make(map[string]*handler)
	}
	e.handlers[topic][id] = sub

	return id, func() { e.Unsubscribe(topic, id) }
}

func (e *EventBus) Unsubscribe(topic string, id string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	handlers, ok := e.handlers[topic]
	if !ok {
		return
	}

	delete(handlers, id)

	if len(handlers) == 0 {
		delete(e.handlers, topic)
	}
}

func (e *EventBus) Publish(ctx context.Context, event Event) {
	e.mu.RLock()
	handlers := e.handlers[event.Topic()]
	e.mu.RUnlock()

	for _, h := range handlers {
		if !h.async {
			next := e.middleware.Wrap(h.base)
			if err := next.Handle(ctx, event); err != nil {
				e.config.ErrorHandler(err)
			}
			continue
		}

		e.wg.Add(1)
		go func() {
			defer e.wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), e.config.AsyncTimeout)
			defer cancel()

			next := e.middleware.Wrap(h.base)
			if err := next.Handle(ctx, event); err != nil {
				e.config.ErrorHandler(err)
			}
		}()
	}
}

func (e *EventBus) Flush(ctx context.Context, queue *EventQueue) {
	for _, event := range queue.Release() {
		e.Publish(ctx, event)
	}
}

func (e *EventBus) Use(middleware ...Middleware) {
	e.middleware.Append(middleware...)
}

func (e *EventBus) Wait() {
	e.wg.Wait()
}

func generateID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x", b)
}
