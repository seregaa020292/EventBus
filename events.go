package eventbus

import "sync"

type Events struct {
	queue []Event
	mu    sync.Mutex
}

func (e *Events) Enqueue(event Event) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.queue = append(e.queue, event)
}

func (e *Events) Release() []Event {
	e.mu.Lock()
	defer e.mu.Unlock()

	events := make([]Event, len(e.queue))
	copy(events, e.queue)

	e.queue = nil
	//e.queue = e.queue[:0]

	return events
}