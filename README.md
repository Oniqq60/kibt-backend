Go-сервис для приема заявок, валидации, сохранения в PostgreSQL и отправки SMTP-уведомлений.  
Стек: `chi/v5` • `pgx/v5` • `validator/v10` • `gomail.v2` • `slog` • Docker

## Переменные окружения
Создайте файл `.env` в папке `kibit-backend/`:

| Переменная | Описание | По умолчанию |
|---|---|---|
| `SERVER_PORT` | Порт HTTP-сервера | `8080` |
| `DB_HOST`, `DB_PORT` | Адрес и порт PostgreSQL | `localhost:5432` |
| `DB_USER`, `DB_PASSWORD`, `DB_NAME` | Учетные данные БД | `postgres` / `password` / `kibt_leads` |
| `SMTP_HOST`, `SMTP_PORT`, `SMTP_USER`, `SMTP_PASS`, `SMTP_FROM` | Настройки SMTP для отправки писем | Gmail конфигурация |
| `RATE_LIMIT_RPS` | Запросов в секунду с одного IP | `10` |
| `RATE_LIMIT_BURST` | Максимальный размер токенов | `20` |
| `LOG_LEVEL` | Уровень логирования (`info`, `debug`, `error`) | `info` |

### Изменение сервиса уведомлений  
В файле 
```bash
internal\service\lead.go строчка 25
```
Вместо sales@kibt.ru вставить почту получателя (админа)

```bash
Так же в .env изменить переменные (SMTP_USER SMTP_PASS SMTP_FROM) на свои 
SMTP_USER - почта от куда отправляються уведомления
SMTP_PASS - app password (для этого нужно включить двухфакторную аунтефикацию и можно будет поставить app password)
SMTP_FROM - изменить почту в ковычках на SMTP_USER
```

## Быстрый старт (локально)

### 1. Клонирование и установка

```bash
cd kibit-backend
go mod download
```

### 2. Настройка окружения
Создайте файл и отредактируйте .env:

```bash
SERVER_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=kibt_leads
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=balandafm@gmail.com
SMTP_PASS=jpcg zhhh zzig mzty
SMTP_FROM="KIBT Leads <balandafm@gmail.com>"
RATE_LIMIT_RPS=10 #  Пропуск запросов в секунду от одного ip
RATE_LIMIT_BURST=20 # Максимальный размер ведра с токенами
LOG_LEVEL=info
```
Или перименовать .env.example и .env

3. Запуск сервера
go run cmd/server/main.go

## Быстрый старт (Контейнер)

### 1. Настройка окружения
Создайте файл и отредактируйте .env:

```bash
SERVER_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=kibt_leads
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=balandafm@gmail.com
SMTP_PASS=jpcg zhhh zzig mzty
SMTP_FROM="KIBT Leads <balandafm@gmail.com>"
RATE_LIMIT_RPS=10 #  Пропуск запросов в секунду от одного ip
RATE_LIMIT_BURST=20 # Максимальный размер ведра с токенами
LOG_LEVEL=info
```

### 2. Запуск через Docker Compose

```bash
docker-compose up -d
```

## API Reference

### 1. Health Check
GET /healthz
Проверка работоспособности сервиса.
Response:

``` json
{
  "status": "ok",
  "timestamp": "2026-05-12T20:00:00Z"
}
```

Status Codes:
200 OK - сервис работает

### 2. Создание заявки
POST /api/v1/leads
Создает новую заявку, валидирует данные, сохраняет в БД и отправляет email.
Request Headers:

```bash
Content-Type: application/json
```

```json
{
  "name": "Иван Петров",
  "email": "ivan@company.ru",
  "phone": "+79991234567",
  "company": "ООО Пример",
  "app_type": "mobile",
  "message": "Нужно разработать мобильное приложение для доставки"
}
```

Response 201 Created:

```json
{
  "id": 1,
  "message": "Заявка успешно создана",
  "data": {
    "id": 1,
    "name": "Иван Петров",
    "email": "ivan@company.ru",
    "phone": "+79991234567",
    "company": "ООО Пример",
    "app_type": "mobile",
    "message": "Нужно разработать мобильное приложение для доставки",
    "created_at": "2026-05-12T20:00:00Z"
  }
}
```

Response 400 Bad Request:

```json
{
  "error": "Validation failed",
  "details": [
    "name: minimum 2 characters required",
    "email: invalid email format"
  ]
}
```

Response 429 Too Many Requests:

```json
{
  "error": "Rate limit exceeded. Try again later"
}
```

### 3. Status Codes:
```bash
201 Created - заявка создана
400 Bad Request - ошибка валидации
429 Too Many Requests - превышен лимит запросов
500 Internal Server Error - ошибка сервера
```
