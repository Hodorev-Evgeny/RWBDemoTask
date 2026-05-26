# RWBDemoTask

Сервис для расчёта актуального Top-N поисковых запросов за последние 5 минут.

Сервис читает поток поисковых событий из брокера сообщений, агрегирует запросы и отдаёт актуальный топ через HTTP API.

## Что реализовано

- Консьюмер поисковых событий из брокера сообщений.
- HTTP API для получения Top-N запросов.
- Динамический stop-list.
- Хранение stop-list в Redis.
- Anti-fraud защита от повторной накрутки одинаковых запросов.
- Нагрузочные тесты для HTTP API и ingest-нагрузки.

---

## Запуск

### 1. Клонирование проекта

```bash
git clone git@github.com:Hodorev-Evgeny/RWBDemoTask.git
cd RWBDemoTask
```

### 2. Создание `.env`

```bash
cp .env.example .env
```

### 3. Заполнение `.env`

Пример данных, с которыми проект запустится:

```text
HTTP_ADDR=:8080
HTTP_SHUTDOWN_TIMEOUT=30s

REDIS_ADDR=redis-db:6379
REDIS_PASSWORD=admin
REDIS_DATABASE=0

NATS_URL=nats://nats:4222
NATS_NAME=trend-service
NATS_MAX_RECONNECTS=-1
NATS_TIMEOUT=5s

LOGGER_LEVEL=DEBUG

TIME_ZONE=UTC
```
Для запуска приложения в Docker:
```text
NATS_URL=nats://nats:4222
REDIS_ADDR=redis-db:6379
```
Для локального запуска приложения с хоста:
```text
NATS_URL=nats://localhost:4222
REDIS_ADDR=localhost:6379
```
### 4. Поднять окружение

```bash
make env-up
```

### 5. Запустить приложение локально

```bash
make app-run
```

### 6. Запустить приложение в Docker

```bash
make deploy-run
```

### 7. Остановить приложение

```bash
make deploy-stop
```

### 8. Остановить окружение

```bash
make env-down
```

---

## API

### Получить Top-N запросов

```bash
curl "http://localhost:8080/toplist?limit=10"
```

Пример ответа:

```json
[
  {
    "query": "iphone 15",
    "count": 2
  },
  {
    "query": "nike",
    "count": 1
  }
]
```

Параметры:

| Параметр | Тип | Описание |
|---|---|---|
| `limit` | int | Количество элементов в топе |

Примеры:

```bash
curl "http://localhost:8080/toplist?limit=1"
```

```bash
curl "http://localhost:8080/toplist?limit=5"
```

```bash
curl "http://localhost:8080/toplist?limit=10"
```

Ошибочные значения:

```bash
curl -i "http://localhost:8080/toplist?limit=0"
```

```bash
curl -i "http://localhost:8080/toplist?limit=-1"
```

```bash
curl -i "http://localhost:8080/toplist?limit=abc"
```

---

## Stop-list API

### Получить stop-list

```bash
curl "http://localhost:8080/stoplist"
```

Пример ответа:

```json
{
  "list": ["casino"]
}
```

### Добавить слово в stop-list

```bash
curl -i -X POST "http://localhost:8080/stoplist/casino"
```

Пример ответа:

```text
HTTP/1.1 201 Created

"casino"
```

### Удалить слово из stop-list

```bash
curl -i -X DELETE "http://localhost:8080/stoplist/casino"
```

Пример ответа:

```text
HTTP/1.1 204 No Content
```

---

## Контракт события из брокера

Сервис читает события из subject:

```text
search.events
```

Формат события:

```json
{
  "query": "iphone 15",
  "user_id": 456,
  "session_id": "session-789",
  "time_event": "2026-05-25T18:45:00Z"
}
```

Описание полей:

| Поле | Тип | Описание |
|---|---|---|
| `query` | string | Поисковый запрос |
| `user_id` | int64 | Идентификатор пользователя |
| `session_id` | string | Идентификатор сессии |
| `time_event` | string | Время события в RFC3339 |

### Почему нужны эти поля

`query` нужен для формирования топа поисковых запросов.

`user_id` и `session_id` нужны для предотвращения накрутки. По этим полям сервис понимает, что один и тот же пользователь или одна и та же сессия пытается повторно отправить одинаковый запрос.

`time_event` нужен для того чтобы при необходимости можно было перейти на расчёт окна по времени события.

---

## Отправка события в брокер
```bash
docker compose exec nats-box nats pub search.events '{
  "query": "iphone 15",
  "user_id": 456,
  "session_id": "session-789",
  "time_event": "2026-05-25T18:45:00Z"
}'
```

Проверить результат:

```bash
curl "http://localhost:8080/toplist?limit=10"
```

Ожидаемый пример ответа:

```json
[
  {
    "query": "iphone 15",
    "count": 1
  }
]
```

---

## Проверка работы

### 1. Проверить пустой top-list

```bash
curl "http://localhost:8080/toplist?limit=10"
```

При чистом запуске ожидается:

```json
[]
```

### 2. Отправить тестовое событие

```bash
docker compose exec nats-box nats pub search.events '{
  "query": "iphone 15",
  "user_id": 456,
  "session_id": "session-789",
  "time_event": "2026-05-25T18:45:00Z"
}'
```

### 3. Проверить top-list

```bash
curl "http://localhost:8080/toplist?limit=10"
```

### 4. Проверить stop-list

Добавить слово:

```bash
curl -i -X POST "http://localhost:8080/stoplist/casino"
```

Проверить список:

```bash
curl "http://localhost:8080/stoplist"
```

Удалить слово:

```bash
curl -i -X DELETE "http://localhost:8080/stoplist/casino"
```

---

## Архитектура

Общий поток данных:

```text
broker events
    ↓
NATS consumer
    ↓
anti-fraud check
    ↓
in-memory storage
    ↓
cached top
    ↓
HTTP API
```

Основные части проекта:

```text
cmd/mainapp
  Точка входа приложения.

cmd/loadgen
  Утилита для генерации тестовой нагрузки.

internal/features/ingest
  Чтение событий из брокера сообщений.

internal/features/toplist
  Получение Top-N поисковых запросов.

internal/features/stoplist
  Управление stop-list.

internal/core/storage
  Хранение событий, счётчиков и cached top.

internal/core/repository/redis
  Работа с Redis для anti-fraud и stop-list.
```

### Почему выбрана такая архитектура

#### Top-list

Top-list хранится в памяти Go-приложения. Это позволяет не обращаться к базе данных при каждом HTTP-запросе и не пересчитывать рейтинг каждый раз заново.

Сервис хранит события в скользящем окне 5 минут. Фоновая горутина периодически удаляет устаревшие события и пересобирает `cached top`.

При добавлении события также выполняется проверка на накрутку. Для этого сервис обращается в Redis. Если такой запрос уже есть в Redis, событие не учитывается. Anti-fraud ключ живёт 30 секунд.

#### Stop-list

Stop-list хранится в памяти Go-приложения и дополнительно сохраняется в Redis. Это позволяет перезапускать сервис без потери stop-list.

Хранение stop-list в памяти ускоряет проверку, потому что сервису не нужно обращаться в Redis при каждой пересборке top-list.

#### Хранение в памяти Go

Такое решение выбрано, потому что сервис часто обращается к структурам top-list и stop-list.

Если при каждом запросе обращаться к внешнему хранилищу, это заметно снизит скорость работы сервиса.

---

## Бизнес-логика

### Top-N за последние 5 минут

Сервис учитывает поисковые события только в рамках временного окна 5 минут.

При запуске storage параллельно запускается горутина, которая обновляет список. Она сравнивает текущее время и время события. Если событие старше 5 минут, оно удаляется из текущего окна.

### Anti-fraud

Для защиты от накрутки используется проверка повторов по комбинации:

```text
user_id + session_id + query
```

При первом событии ключ записывается с TTL. Повторное событие с той же комбинацией не учитывается.

Текущий TTL — 30 секунд. При попытке накрутки сервис проверяет, есть ли такой ключ в Redis. Если ключ уже есть, событие не учитывается.

### Stop-list

Stop-list скрывает нежелательные запросы из выдачи top-list.

Слова можно добавлять и удалять без перезапуска сервиса.

Stop-list влияет только на вывод top-list. Если удалить слово из stop-list, top-list пересоберётся, и удалённое слово снова будет отображаться, если событие ещё находится в 5-минутном окне.

---

## Нагрузочное тестирование

### HTTP benchmark

Запуск:

```bash
make start-http-test
```

Пример результата:

```text
Summary:
  Total:        30.0008 secs
  Slowest:      0.0123 secs
  Fastest:      0.0001 secs
  Average:      0.0030 secs
  Requests/sec: 162574.5133
  
  Total data:   780377760 bytes
  Size/request: 780 bytes

Response time histogram:
  0.000 [1]     |
  0.001 [947887]        |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.003 [46918] |■■
  0.004 [4230]  |
  0.005 [667]   |
  0.006 [166]   |
  0.007 [18]    |
  0.009 [12]    |
  0.010 [57]    |
  0.011 [38]    |
  0.012 [6]     |


Latency distribution:
  10%% in 0.0002 secs
  25%% in 0.0003 secs
  50%% in 0.0005 secs
  75%% in 0.0008 secs
  90%% in 0.0010 secs
  95%% in 0.0013 secs
  99%% in 0.0021 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0000 secs, 0.0000 secs, 0.0084 secs
  DNS-lookup:   0.0000 secs, 0.0000 secs, 0.0082 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0038 secs
  resp wait:    0.0027 secs, 0.0000 secs, 0.0119 secs
  resp read:    0.0002 secs, 0.0000 secs, 0.0057 secs

Status code distribution:
  [200] 1000000 responses
```

---

### Ingest benchmark

Запуск:

```bash
make start-nats-test
```

Пример результата:

```text
sent=3000000 failed=0 duration=7.376728625s messages_per_sec=406684.34
```

---

## Trade-offs

### 1. In-memory storage

Текущий top-list хранится в памяти приложения.

Плюсы:

- быстрое чтение;
- простая реализация;
- нет обращения к внешнему хранилищу на каждый запрос.

Минусы:

- при рестарте текущий top-list сбрасывается.

### 2. Redis для защиты от накрутки

Плюсы:

- TTL автоматически очищает временные ключи;
- простая защита от повторов.

Минусы:

- Redis-запрос на каждое событие ограничивает скорость ingest-обработки.

### 3. Stop-list

Плюсы:

- можно менять без перезапуска;
- данные stop-list сохраняются между рестартами.

Минусы:

- stop-list скрывает запрос из выдачи, но не удаляет его из уже накопленного окна.

---
