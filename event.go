package eventbus

import "sync"

type Event interface {
	Topic() string
}

type EventQueue struct {
	queue []Event
	mu    sync.Mutex
}

func (e *EventQueue) Enqueue(event Event) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.queue = append(e.queue, event)
}

func (e *EventQueue) Release() []Event {
	e.mu.Lock()
	defer e.mu.Unlock()

	events := e.queue
	e.queue = nil
	return events
}
