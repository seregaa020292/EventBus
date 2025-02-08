# EventBus

Шина событий на с поддержкой middleware и асинхронной обработкой.

## Особенности

- 🚀 Асинхронная обработка событий
- 🛠️ Поддержка цепочек middleware
- 📦 Гибкая конфигурация
- 🛡️ Обработка ошибок
- 🕒 Таймауты для асинхронных операций

## Быстрый старт

```go
package main

import (
    "context"

    "ppr.gitlab.yandexcloud.net/ecosystem/fines/service/pkg/eventbus"
)

type UserCreatedEvent struct{}

func (e UserCreatedEvent) Topic() string {
    return "user.created"
}

func main() {
    bus := eventbus.New()

    // Подписка на событие
    id, _ := bus.Subscribe("user.created", eventbus.HandlerFunc(func(ctx context.Context, e eventbus.Event) error {
        // Обработка события
        return nil
    }), eventbus.WithAsync())

    // Публикация события
    bus.Publish(context.Background(), &UserCreatedEvent{})

    bus.Wait()
    bus.Unsubscribe("user.created", id)
}

```

## Конфигурация

```go
package main

import (
    "time"

    "ppr.gitlab.yandexcloud.net/ecosystem/fines/service/pkg/eventbus"
)

func main() {
    bus := eventbus.New(
        // Таймаут для асинхронных обработчиков
        eventbus.WithAsyncTimeout(30*time.Second),
        // Кастомный обработчик ошибок
        eventbus.WithErrorHandler(func(err error) {
            // Кастомная обработка ошибок
        }),
    )
}

```

## Middleware

```go
package main

import (
    "context"
    "log"

    "ppr.gitlab.yandexcloud.net/ecosystem/fines/service/pkg/eventbus"
)

func main() {
    bus := eventbus.New()

    // Пример middleware для логирования
    bus.Use(func(next eventbus.Handler) eventbus.Handler {
        return eventbus.HandlerFunc(func(ctx context.Context, e eventbus.Event) error {
            log.Printf("Обработка события: %s", e.Topic())
            return next.Handle(ctx, e)
        })
    })
}

// Порядок выполнения:
// 1. Первый зарег. middleware
// 2. Второй зарег. middleware
// 3. ...
// 4. Основной обработчик

```

## Очередь событий

```go
package main

import (
    "context"

    "ppr.gitlab.yandexcloud.net/ecosystem/fines/service/pkg/eventbus"
)

func main() {
    bus := eventbus.New()

    queue := &eventbus.EventQueue{}

    // Накопление событий
    queue.Enqueue(event1)
    queue.Enqueue(event2)

    // Пакетная публикация
    bus.Flush(context.Background(), queue)
}
```
