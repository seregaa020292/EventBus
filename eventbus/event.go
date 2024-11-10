package eventbus

type EventType uint16
type EventOption func(*Event)

type Event struct {
	typ     EventType
	payload any
	isAsync bool
}

func NewEvent(typ EventType, payload any, options ...EventOption) Event {
	e := Event{
		typ:     typ,
		payload: payload,
		isAsync: false,
	}

	for _, option := range options {
		option(&e)
	}

	return e
}

func WithIsAsync(state bool) EventOption {
	return func(e *Event) {
		e.isAsync = state
	}
}
