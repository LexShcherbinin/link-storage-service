# 🔗 Link Storage Service

Сервис управления ссылками. Позволяет сохранять ссылки, получать их по короткому идентификатору и вести статистику обращений.

---

# 🚀 Функциональность

* Создание короткой ссылки (short URL)
* Получение оригинальной ссылки по shortCode
* Подсчёт количества переходов (visits)
* Кэширование через Redis
* Миграции базы данных
* Graceful shutdown
* Docker-окружение

---

# 🧱 Архитектура

Проект построен по layered architecture:

```
handler → service → repository → database
                 ↘ cache (Redis)
```

### Слои:

* **handler** — HTTP-слой (обработка запросов)
* **service** — бизнес-логика
* **repository** — работа с PostgreSQL
* **cache** — Redis для ускорения чтения

---

# 🛠️ Технологии

* Go
* PostgreSQL
* Redis
* Docker / docker-compose
* golang-migrate

---

# 📦 Структура проекта

```
.
├── cmd/app                # точка входа
├── internal/
│   ├── handler           # HTTP handlers
│   ├── service           # бизнес-логика
│   ├── repository        # работа с БД
│   ├── cache             # Redis
│   ├── model             # модели
│   └── middleware        # middleware
├── migrations            # SQL миграции
├── Dockerfile
├── docker-compose.yml
└── README.md
```

---

# ⚙️ Запуск проекта

## 🐳 Через Docker (рекомендуется)

```bash
docker-compose up --build
```

После запуска:

* API: http://localhost:8080
* PostgreSQL: localhost:5432
* Redis: localhost:6379

---

## 🛑 Остановка

```bash
Ctrl + C
```

или

```bash
docker-compose down
```

---

# 🔄 Миграции

Миграции применяются автоматически при старте приложения.

Файлы находятся в:

```
migrations/
  000001_create_links_table.up.sql
  000001_create_links_table.down.sql
```

---

# 📡 API

## ➕ Создать короткую ссылку

```
POST /links
```

### Request:

```json
{
  "url": "https://example.com"
}
```

### Response:

```json
"b"
```

### Описание:

* Принимает оригинальный URL
* Создаёт запись в БД
* Генерирует shortCode (base62)
* Возвращает короткий код

---

## 📋 Получить список ссылок

```
GET /links
```

### Response:

```json
[
  {
    "short_code": "b",
    "original_url": "https://example.com",
    "visits": 3
  }
]
```

### Описание:

* Возвращает список всех сохранённых ссылок
* Используется для просмотра данных в системе

---

## 🔍 Получить оригинальный URL

```
GET /links/{code}
```

### Response:

```json
{
  "short_code": "b",
  "original_url": "https://example.com",
  "visits": 3
}
```

### Описание:

* Принимает shortCode
* Сначала проверяет Redis (кэш)
* Если нет — идёт в PostgreSQL
* Увеличивает счётчик visits
* Возвращает данные ссылки

---

## 📊 Получить статистику ссылки

```
GET /links/{short_code}/stats
```

### Response:

```json
{
  "short_code": "b",
  "visits": 3,
  "created_at": "2026-01-01T12:00:00Z"
}
```

### Описание:

* Возвращает статистику по ссылке
* Не использует кэш (берёт актуальные данные из БД)
* Показывает количество переходов и дату создания

---

## ❌ Удалить ссылку

```
DELETE /links/{code}
```

### Response:

```
204 No Content
```

### Описание:

* Удаляет ссылку из базы данных
* Удаляет запись из Redis-кэша
* Используется для очистки данных

---

# 🧠 Как работает генерация shortCode

* используется auto-increment ID из PostgreSQL
* кодируется в base62
* гарантирует уникальность без дополнительных проверок

---

# ⚡ Кэширование

* Redis хранит: `short_code → original_url`
* уменьшает нагрузку на БД
* БД остаётся источником истины

---

# 📈 Подсчёт переходов

* увеличивается при каждом запросе ссылки
* выполняется атомарно в PostgreSQL:

```sql
UPDATE links SET visits = visits + 1
```

---

# 🐳 Docker

Сервис запускается в 3 контейнерах:

* app
* postgres
* redis

Контейнеры взаимодействуют через внутреннюю сеть Docker:

```
app → postgres:5432
app → redis:6379
```

---

# 🔐 Конфигурация

Через environment variables:

```
DB_URL=postgres://postgres:postgres@postgres:5432/links?sslmode=disable
REDIS_ADDR=redis:6379
```

---

# 🧪 Проверка

Создать ссылку:

```bash
curl -X POST http://localhost:8080/links \
  -H "Content-Type: application/json" \
  -d '{"url":"https://example.com"}'
```

Получить ссылку:

```bash
curl http://localhost:8080/links/{code}
```

---
