# 🔗 Link Storage Service

Сервис для сокращения ссылок с возможностью хранения и получения оригинального URL, а также сбора статистики переходов.

---

## 🚀 Функциональность

* Создание короткой ссылки
* Получение оригинального URL по short_code
* Подсчёт количества переходов
* Получение списка ссылок с пагинацией
* Удаление ссылок
* Просмотр статистики

---

## 🛠️ Технологии

* Go (Golang)
* PostgreSQL
* golang-migrate (миграции)
* Docker (для БД)

---

## 📁 Структура проекта

```
.
├── cmd/app              # Точка входа (main.go)
├── internal/
│   ├── handler/         # HTTP handlers
│   ├── service/         # Бизнес-логика
│   ├── repository/      # Работа с БД
│   ├── model/           # Модели
│   ├── cache/           # Кеш
│   └── config/          # Конфигурация
├── migrations/          # SQL миграции
├── go.mod
└── README.md
```

---

## ⚙️ Запуск проекта

### 1. Запуск PostgreSQL через Docker

```bash
docker run --name link-db \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=links \
  -p 5432:5432 \
  -d postgres
```

---

### 2. Установка migrate

#### macOS

```bash
brew install golang-migrate
```

#### Linux (пример)

```bash
curl -L https://github.com/golang-migrate/migrate/releases/latest/download/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/
```

---

### 3. Применение миграций

```bash
migrate -path ./migrations \
  -database "postgres://postgres:postgres@localhost:5432/links?sslmode=disable" \
  up
```

---

### 4. Настройка переменных окружения

```bash
export DB_URL=postgres://postgres:postgres@localhost:5432/links?sslmode=disable
export PORT=8080
```

---

### 5. Запуск сервиса

```bash
go run cmd/app/main.go
```

---

## 📡 API

---

### ➕ Создание короткой ссылки

**POST /links**

#### Request

```json
{
  "url": "https://example.com/some/very/long/url"
}
```

#### Response

```json
{
  "short_code": "abc123"
}
```

---

### 🔍 Получение оригинальной ссылки

**GET /links/{short_code}**

#### Response

```json
{
  "url": "https://example.com",
  "visits": 10
}
```

---

### 📄 Получение списка ссылок

**GET /links?limit=10&offset=0**

---

### ❌ Удаление ссылки

**DELETE /links/{short_code}**

---

### 📊 Статистика

**GET /links/{short_code}/stats**

#### Response

```json
{
  "short_code": "abc123",
  "url": "https://example.com",
  "visits": 10,
  "created_at": "2026-01-01T12:00:00Z"
}
```

---

## ⚡ Особенности реализации

* Генерация уникального `short_code`
* Потокобезопасное увеличение счётчика переходов
* Кеширование часто запрашиваемых ссылок
* Чистая архитектура (handler → service → repository)
* Конфигурация через переменные окружения

---

## 🧪 Пример использования (curl)

```bash
curl -X POST http://localhost:8080/links \
  -H "Content-Type: application/json" \
  -d '{"url":"https://example.com"}'
```

---

## 📌 TODO (дополнительно)

* Добавить Redis для кеша
* Логи (zap/logrus)
* Graceful shutdown
* Docker Compose для полного окружения

---

## 👨‍💻 Автор

Тестовое задание
