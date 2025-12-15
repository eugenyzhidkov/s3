# S3 Files Manager

**HTTP-сервис для работы с S3-совместимым хранилищем** — демонстрация загрузки, сжатия и управления файлами в MinIO (локальный S3).
Сервис предоставляет веб-интерфейс для просмотра, загрузки и перемещения файлов в S3.

## Что реализовано

- HTTP-сервер на порту из env (`HTTP_PORT`, дефолт 8080)  
- Полный **graceful shutdown** (обработка сигналов, таймаут через `SHUTDOWN_TUMEOUT`)  
- Конфигурация **полностью через env** с использованием `caarlos0/env/v11`
- Docker Compose с MinIO (healthcheck + гибкие порты) 
- Веб-интерфейс: список файлов в бакете `my-bucket`
- Загрузка JPEG из папки `static/` с сжатием (quality 75%, ресайз до 1920px)
- Перемещение файлов в виртуальную папку `archive/` (copy + delete)
- Чистая архитектура: `cmd/`, `internal/`
- Логирование через Zap  

## Структура проекта
s3/
├── cmd/main.go                        # точка входа
├── internal/
│   ├── app/app.go                     # оркестратор запуска и graceful shutdown
│   ├── config/config.go               # конфигурация (env)
│   ├── handler/                       # HTTP-хендлеры и UI
│   │   ├── index_handler.go
│   │   ├── upload_handler.go
│   │   └── move_handler.go
│   ├── image/image.go              # сжатие изображений
│   ├── infrastructure/logger/         # Zap-логгер
│   └── s3/                            # S3-клиент (AWS SDK v2)
│       ├── s3.go
│       ├── list.go
│       ├── upload.go
│       └── move.go
├── static/                            # исходные картинки для загрузки
├── templates/index.html               # веб-интерфейс
├── docker-compose.yml
├── .env
├── go.mod / go.sum
└── Makefile

## Переменные окружения

| Переменная              | Описание                          | По умолчанию    |
|-------------------------|-----------------------------------|-----------------|
| `HTTP_PORT`             | Порт gRPC-сервера                 | `8080`          |
| `S3_ENDPOINT`           | Endpoint MinIO/S3                 | `localhost:9000`|
| `S3_ACCESS_KEY`         | Access key                        | `minioadmin`    |
| `S3_SECRET_KEY`         | Secret key                        | `minioadmin`    |
| `S3_BUCKET`             | Имя бакета                        | `my_bucket`     |
| `S3_USE_SSL`            | Использовать HTTPS                | `false`         |
| `SHUTDOWN_TIMEOUT`      | Таймаут graceful shutdown         | `10s`           |
| `LOG_LEVEL`             | Уровень логирования               | `info`          |

## Запуск сервиса

### 1. Запустить MinIO
make minio-up

### 2. Создать бакет my-bucket в http://localhost:9001

### 3. Положить JPEG-файлы в папку static/

### 4. Запуск сервиса
make run

### 5. Запуск всего окружения (Docker Compose)
docker compose up --build

## Ожидаемый вывод при старте

S3-клиент создан endpoint=localhost:9000 bucket=my-bucket
подключение к S3 проверено успешно
HTTP-сервер запущен port=8080

## HTTP эндпоинты
**GET /** — главная страница с списком файлов в бакете
**POST /upload** — загрузка и сжатие всех JPEG из `static`/` в S3
**POST /move** — перемещение всех файлов в виртуальную папку `archive/`

Сервис полностью готов к дальнейшему развитию (добавление авторизации, поддержка других форматов, интеграция с реальным S3).