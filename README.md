# todo

## Архитектура

Проект построен по луковичной архитектуре и разбит на три микросервиса:

- **auth-service** (порт 8081) — регистрация, логин, валидация JWT токенов. БД: PostgreSQL
- **tasks-service** (порт 8082) — управление задачами и заметками. БД: PostgreSQL + MongoDB + Redis
- **notifier-service** — читает события из Kafka, создаёт уведомления при смене статуса задачи

## Схема взаимодействия
Client → auth-service (register/login)

Client → tasks-service (tasks/notes) → auth-service (validate token)

tasks-service → Kafka → notifier-service → MongoDB (notifications)

### Стек
| Компонент | Технология |
|-----------|-----------|
| Язык | Go |
| HTTP роутер | Chi |
| PostgreSQL | sqlx + lib/pq |
| MongoDB | mongo-driver |
| Кеш | Redis |
| Брокер | Kafka |
| Миграции | Goose |

## Как поднять локально

1. Склонировать репозиторий
2. Скопировать `.env.example` в `.env` и заполнить переменные
3. Запустить:

```bash
docker compose up --build
```

Сервисы будут доступны:
- auth-service: http://localhost:8081
- tasks-service: http://localhost:8082
- Kafka UI: http://localhost:8090

## Переменные окружения

| Переменная | Описание | Пример |
|------------|----------|--------|
| APP_PORT | Порт приложения | 8080 |
| JWT_SECRET | Секрет для JWT | your-secret |
| JWT_TTL | Время жизни токена | 24h |
| AUTH_POSTGRES_HOST | Хост PostgreSQL для auth | postgres-auth |
| TASKS_POSTGRES_HOST | Хост PostgreSQL для tasks | postgres-tasks |
| MONGO_HOST | Хост MongoDB | mongo |
| REDIS_HOST | Хост Redis | redis |
| KAFKA_BROKERS | Адрес Kafka | kafka:9092 |
| AUTH_SERVICE_URL | URL auth-service | http://auth-service:8081 |

## Примеры запросов

### Регистрация
```bash
curl -X POST http://localhost:8081/api/register \
  -H "Content-Type: application/json" \
  -d '{"name": "User", "email": "user@test.com", "password": "password123"}'
```

### Логин
```bash
curl -X POST http://localhost:8081/api/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user@test.com", "password": "password123"}'
```

### Создать задачу
```bash
curl -X POST http://localhost:8082/api/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"title": "My task", "description": "Description"}'
```

### Получить список задач
```bash
curl http://localhost:8082/api/tasks \
  -H "Authorization: Bearer <token>"
```

### Обновить статус задачи
```bash
curl -X PUT http://localhost:8082/api/tasks/{id} \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"title": "My task", "description": "Description", "status": "in_progress"}'
```

### Добавить заметку
```bash
curl -X POST http://localhost:8082/api/tasks/{id}/notes \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"text": "My note"}'
```

## Выбор библиотек

- **Chi** — лёгкий идиоматичный роутер, близкий к стандартной библиотеке
- **sqlx** — тонкая обёртка над database/sql, полный контроль над SQL запросами
- **mongo-driver** — официальный драйвер MongoDB
- **go-redis** — стандарт де-факто для Redis в Go
- **segmentio/kafka-go** — простой идиоматичный клиент Kafka без CGO зависимостей
- **golang-jwt/jwt** — популярная библиотека для JWT
- **goose** — простые SQL миграции
- **slog** — стандартная библиотека Go для структурированных логов

## Что можно улучшить

- Добавить больше тестов
- Добавить пагинацию для списка задач
- Улучшить readyz — реальная проверка зависимостей
- Добавить трейсинг
- Вынести конфигурацию goose в отдельную команду
