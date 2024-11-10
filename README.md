## Usage

```go
package main

const (
	EventTypeCreated eventbus.EventType = iota
	EventTypeUpdated
)

type ListenHandle1 struct {
	V int
}

func (h ListenHandle1) Handle(ctx context.Context, payload any) {
	fmt.Printf("eventbus handler1 event: %+v, %d\n", payload, h.V)
}

func main() {
	bus := eventbus.New()

	listenHandle1 := ListenHandle1{1}
	listenHandle2 := ListenHandle1{2}

	{
		_, unsubscribe := bus.Subscribe(EventTypeCreated, listenHandle1)
		unsubscribe()
	}
	{
		handlerID, _ := bus.Subscribe(EventTypeCreated, listenHandle2)
		bus.Unsubscribe(EventTypeCreated, handlerID)
	}

	handler2 := eventbus.HandlerFunc(func(ctx context.Context, payload any) {
		switch e := payload.(type) {
		case eventCreated:
			fmt.Printf("eventbus handler2 event: %+v, %v\n", e, e.ID)
		case eventUpdated:
			fmt.Printf("eventbus handler2 event: %+v, %v\n", e, e.ID)
		}
	})
	bus.Subscribe(EventTypeCreated, handler2)
	bus.Subscribe(EventTypeUpdated, handler2)

	bus.Publish(context.Background(), eventbus.NewEvent(
		EventTypeCreated,
		eventCreated{"Foo eventCreated"},
	))

	var events = eventbus.Events{}

	events.Enqueue(eventbus.NewEvent(
		EventTypeCreated,
		eventCreated{"Foo eventCreated12"},
		eventbus.WithIsAsync(true),
	))
	events.Enqueue(eventbus.NewEvent(
		EventTypeUpdated,
		eventUpdated{"Foo eventUpdated122"},
	))

	bus.Flush(context.Background(), &events)

	time.Sleep(200 * time.Millisecond)
}
```
