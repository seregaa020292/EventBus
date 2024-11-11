package eventbus_test

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"sanbox/eventbus"
)

type testHandler struct {
	mu      sync.Mutex
	calls   []eventbus.Event
	handler func(ctx context.Context, event eventbus.Event)
}

func (h *testHandler) Handle(ctx context.Context, event eventbus.Event) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.calls = append(h.calls, event)
	if h.handler != nil {
		h.handler(ctx, event)
	}
}

type testEvent string

func (e testEvent) Name() eventbus.EventName { return 1 }
func (e testEvent) IsAsync() bool            { return false }

type testEventAsync string

func (e testEventAsync) Name() eventbus.EventName { return 1 }
func (e testEventAsync) IsAsync() bool            { return true }

func TestEventBus(t *testing.T) {
	bus := eventbus.New()
	ctx := context.Background()

	// Создаем обработчик
	var calls []eventbus.Event
	handler := &testHandler{
		handler: func(ctx context.Context, event eventbus.Event) {
			calls = append(calls, event)
		},
	}

	// Подписываемся на событие
	event := testEvent("test payload")
	_, unsubscribe := bus.Subscribe(event.Name(), handler)

	// Проверяем, что подписка работает
	bus.Publish(ctx, event)

	// Убедимся, что обработчик был вызван
	assert.Len(t, calls, 1)
	assert.Equal(t, testEvent("test payload"), calls[0])

	// Отписываемся и проверяем, что обработчик больше не вызывается
	unsubscribe()
	bus.Publish(ctx, testEvent("another payload"))

	// Убедимся, что обработчик не был вызван
	assert.Len(t, calls, 1)
}

func TestFlush(t *testing.T) {
	bus := eventbus.New()
	ctx := context.Background()
	events := &eventbus.Events{}

	// Создаем обработчик
	var calls []any
	handler := &testHandler{
		handler: func(ctx context.Context, event eventbus.Event) {
			calls = append(calls, event)
		},
	}

	// Подписываемся на событие
	bus.Subscribe(eventbus.EventName(1), handler)

	// Добавляем события в очередь
	events.Enqueue(testEvent("flush payload 1"))
	events.Enqueue(testEvent("flush payload 2"))

	// Вызываем Flush
	bus.Flush(ctx, events)

	// Проверяем, что оба события были обработаны
	assert.Len(t, calls, 2)
	assert.Equal(t, testEvent("flush payload 1"), calls[0])
	assert.Equal(t, testEvent("flush payload 2"), calls[1])
}

func TestAsyncPublish(t *testing.T) {
	bus := eventbus.New()
	ctx := context.Background()
	wg := &sync.WaitGroup{}

	wg.Add(1)

	// Создаем обработчик
	var calls []any
	handler := &testHandler{
		handler: func(ctx context.Context, event eventbus.Event) {
			defer wg.Done()
			calls = append(calls, event)
		},
	}

	// Подписываемся на событие
	bus.Subscribe(eventbus.EventName(1), handler)

	// Публикуем асинхронное событие
	bus.Publish(ctx, testEventAsync("async payload"))

	// Даем время для обработки
	// Ждем немного, чтобы асинхронное событие было обработано
	wg.Wait()

	// Проверяем, что событие было обработано
	assert.Len(t, calls, 1)
	assert.Equal(t, testEventAsync("async payload"), calls[0])
}

func TestUnsubscribe(t *testing.T) {
	bus := eventbus.New()
	ctx := context.Background()

	// Создаем обработчик
	var calls []any
	handler := &testHandler{
		handler: func(ctx context.Context, event eventbus.Event) {
			calls = append(calls, event)
		},
	}

	// Подписываемся на событие
	handlerID, unsubscribe := bus.Subscribe(eventbus.EventName(1), handler)

	assert.Equal(t, uint16(1), handlerID)

	// Публикуем событие
	bus.Publish(ctx, testEvent("payload before unsubscribe"))

	// Убедимся, что обработчик был вызван
	assert.Len(t, calls, 1)

	// Отписываемся
	unsubscribe()

	// Публикуем еще одно событие
	bus.Publish(ctx, testEvent("payload after unsubscribe"))

	// Убедимся, что обработчик не был вызван
	assert.Len(t, calls, 1)
}
