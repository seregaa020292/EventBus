package eventbus_test

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"sanbox/eventbus"
)

type testHandler struct {
	calls []eventbus.Event
	mu    sync.Mutex
}

func (h *testHandler) Handle(ctx context.Context, event eventbus.Event) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.calls = append(h.calls, event)
	return nil
}

type testEvent string

func (e testEvent) Topic() string { return "test.event" }

func TestEventBus(t *testing.T) {
	bus := eventbus.New()
	ctx := context.Background()

	// Создаем обработчик
	handler := &testHandler{}

	// Подписываемся на событие
	event := testEvent("test payload")
	_, unsubscribe := bus.Subscribe(event.Topic(), handler)

	// Публикуем событие
	bus.Publish(ctx, event)

	// Убедимся, что обработчик был вызван
	assert.Len(t, handler.calls, 1)
	assert.Equal(t, testEvent("test payload"), handler.calls[0])

	// Отписываемся и проверяем, что обработчик больше не вызывается
	unsubscribe()
	bus.Publish(ctx, testEvent("another payload"))

	// Убедимся, что обработчик не был вызван
	assert.Len(t, handler.calls, 1)
}

func TestFlush(t *testing.T) {
	bus := eventbus.New()
	ctx := context.Background()

	// Создаем обработчик
	handler := &testHandler{}

	// Подписываемся на событие
	event := testEvent("test payload")
	bus.Subscribe(event.Topic(), handler)

	// Добавляем события в очередь
	events := &eventbus.EventQueue{}
	events.Enqueue(event)
	events.Enqueue(event)

	// Вызываем Flush
	bus.Flush(ctx, events)

	// Проверяем, что оба события были обработаны
	assert.Len(t, handler.calls, 2)
	assert.Equal(t, testEvent("test payload"), handler.calls[0])
	assert.Equal(t, testEvent("test payload"), handler.calls[1])
}

func TestAsyncPublish(t *testing.T) {
	bus := eventbus.New()
	ctx := context.Background()

	// Создаем обработчик
	handler := &testHandler{}

	// Подписываемся на событие
	event := testEvent("async payload")
	bus.Subscribe(event.Topic(), handler, eventbus.WithHandlerAsync())

	// Публикуем асинхронное событие
	bus.Publish(ctx, event)

	// Даем время для обработки
	bus.Wait()

	// Проверяем, что событие было обработано
	assert.Len(t, handler.calls, 1)
	assert.Equal(t, testEvent("async payload"), handler.calls[0])
}

func TestUnsubscribe(t *testing.T) {
	bus := eventbus.New()
	ctx := context.Background()

	// Создаем обработчик
	handler := &testHandler{}

	// Подписываемся на событие
	event := testEvent("payload unsubscribe")
	_, unsubscribe := bus.Subscribe(event.Topic(), handler)

	// Публикуем событие
	bus.Publish(ctx, event)

	// Убедимся, что обработчик был вызван
	assert.Len(t, handler.calls, 1)

	// Отписываемся
	unsubscribe()

	// Публикуем еще одно событие
	bus.Publish(ctx, event)

	// Убедимся, что обработчик не был вызван
	assert.Len(t, handler.calls, 1)
}
