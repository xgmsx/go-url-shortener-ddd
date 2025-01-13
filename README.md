# go-url-shortener-ddd

[![Lint status](https://github.com/xgmsx/go-url-shortener-ddd/actions/workflows/golangci-lint.yml/badge.svg?branch=main)](https://github.com/xgmsx/go-url-shortener-ddd/actions/workflows/golangci-lint.yml)
[![Test status](https://github.com/xgmsx/go-url-shortener-ddd/actions/workflows/coverage.yml/badge.svg?branch=main)](https://github.com/xgmsx/go-url-shortener-ddd/actions/workflows/coverage.yml)
[![coverage_report](https://img.shields.io/badge/coverage_report-30.8%25-yellow)](https://xgmsx.github.io/go-url-shortener-ddd)

**go-url-shortener-ddd** — это сокращатель ссылок на Go, реализованный с использованием архитектурного подхода Domain-Driven Design (DDD).

Структура:
* [internal](internal) содержит пакеты относящиеся к проекту - app, config, shortener.
* [pkg](pkg) содержит универсальные пакеты, которые можно переиспользовать - logger, http-сервер, клиенты к postgres, redis, kafka и т.д. 

Table of Contents:
- [Features](#features)
- [Installation](#installation)
- [Quick start](#quick-start)
- [Usage](#usage)
- [Metrics](#metrics)
- [Traces](#traces)
- [Environment variables](#environment-variables)

## Features

- **Поддержка протоколов**: HTTP и gRPC интерфейсы для взаимодействия с сервисом.
- **Интеграция с Kafka**: Отправка и получение сообщений в kafka для взаимодействия с сервисом.
- **Хранение данных**: Данные о созданных ссылках хранятся в Postgres SQL.
- **Кэширование**: Данные о созданных и запрашиваемых ссылках кешируются в Redis для снижения нагрузки на БД.  
- **Трассировки (observability)**: Данные трейсов входящих запросов отправляются в Jaeger для анализа производительности распределенных систем.
- **Метрики (observability)**: Данные метрик входящих запросов отправляются в Prometheus для анализа производительности сервиса в Grafana.
- **Логи (observability)**: Информация об ошибках передается в Sentry. Работа сервиса логируется в формате JSON.

## Installation

<div class="termy">

```console
git clone https://github.com/xgmsx/go-url-shortener-ddd
cd go-url-shortener-ddd
cp ./configs/.env_example ./configs/.env
cd ./configs/.env_localhost_example ./configs/.env_localhost
```

</div>

## Quick start

#### Запуск сервиса с зависимостями (postgres, redis, kafka, jaeger, prometheus, grafana) в Docker:

<div class="termy">

```console
$ docker compose up -d --build

 ✔ Network go-url-shortener-ddd_default         Created                                                                                                                                                                       0.1s 
 ✔ Container go-url-shortener-ddd-app-1         Started                                                                                                                                                                       0.4s 
 ✔ Container go-url-shortener-ddd-postgres-1    Started                                                                                                                                                                       0.2s 
 ✔ Container go-url-shortener-ddd-redis-1       Started                                                                                                                                                                       0.2s 
 ✔ Container go-url-shortener-ddd-redis-ui-1    Started                                                                                                                                                                       0.2s 
 ✔ Container go-url-shortener-ddd-kafka-1       Started                                                                                                                                                                       0.2s 
 ✔ Container go-url-shortener-ddd-kafka-init-1  Started   
 ✔ Container go-url-shortener-ddd-kafka-ui-1    Started                                                                                                                                                                       0.3s 
```

</div>

#### Локальный запуск сервиса в Linux и MacOS:

<div class="termy">

```console
$ export $(grep -v '^#' ./configs/.env_localhost | xargs) && go run ./cmd/app

2025/01/01 12:00:00 maxprocs: Leaving GOMAXPROCS=8: CPU quota undefined
12:00:00 INF Logger initialized
12:00:00 INF App starting...
12:00:00 INF Sentry is disabled
12:00:00 INF Tracer is disabled
12:00:00 INF HTTP server started on port: 8000
12:00:00 INF gRPC server started on port: 50051
12:00:00 INF Kafka consumer started
12:00:00 INF App started
```

</div>

#### Локальный запуск сервиса в Windows Git Bash:

```console
$ env $(grep -v '^#' ./configs/.env_localhost | xargs) go run ./cmd/app

2025/01/01 12:00:00 maxprocs: Leaving GOMAXPROCS=8: CPU quota undefined
12:00:00 INF Logger initialized
12:00:00 INF App starting...
12:00:00 INF Sentry is disabled
12:00:00 INF Tracer is disabled
12:00:00 INF HTTP server started on port: 8000
12:00:00 INF gRPC server started on port: 50051
12:00:00 INF Kafka consumer started
12:00:00 INF App started
```

## Usage

#### Migrations

**Warning**: При первом запуске сервиса нужно выполнить команды  `make migrate-install` и `make migrate-up`:

<div class="termy">

```console
$ make migrate-up

migrate -database "postgres://login:pass@localhost:5432/app-db?sslmode=disable" -path "./migrations" up
20241212085257/u init (19.9317ms)
```

</div>

#### HTTP-запросы

**Note**: Выполнять запросы можно в веб-интерфейсе http://localhost:8000/swagger

Установка утилиты `curl`:
```shell
# Linux
$ sudo apt-get install curl

# MacOS
$ brew install curl

# Windows
$ scoop install curl
```

Создание новой короткой ссылки:
```shell
curl -X 'POST' \
  'http://localhost:8000/api/shortener/v1/link' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "url": "https://google.com"
}'

# {"url":"https://google.com","alias":"IFIYr0OGRKeqF9jPUIbwww","expired_at":"2025-01-02T12:00:00.000000000Z"}
```

Получение полной ссылки:
```shell
curl -X 'GET' \
  'http://localhost:8000/api/shortener/v1/link/IFIYr0OGRKeqF9jPUIbwww' \
  -H 'accept: application/json'

# {"url": "https://google.com", "alias": "IFIYr0OGRKeqF9jPUIbwww", "expired_at": "2025-01-02T12:00:00.000000000Z"}
```

Переход по короткой ссылке:
```shell
# Linux and MacOS
open http://localhost:8000/api/shortener/v1/link/IFIYr0OGRKeqF9jPUIbwww/redirect

# Windows:
explorer http://localhost:8000/api/shortener/v1/link/IFIYr0OGRKeqF9jPUIbwww/redirect
```

#### gRPC-запросы

Установка утилиты `grpcurl`:
```shell
# Linux
$ sudo apt-get install grpcurl

# MacOS
$ brew install grpcurl

# Windows
$ scoop install grpcurl
```

Создание новой короткой ссылки:
```shell
$ grpcurl -d '{"url": "https://google.com"}' -plaintext localhost:50051 shortener_v1.Shortener/CreateLink

# {"url": "https://google.com", "alias": "IFIYr0OGRKeqF9jPUIbwww", "expired_at": "2025-01-02T12:00:00.000000000Z"}
```

Получение полной ссылки:
```shell
$ grpcurl -d '{"alias": "IFIYr0OGRKeqF9jPUIbwww"}' -plaintext localhost:50051 shortener_v1.Shortener/GetLink

# {"url": "https://google.com", "alias": "IFIYr0OGRKeqF9jPUIbwww", "expired_at": "2025-01-02T12:00:00.000000000Z"}
```

#### Redis UI

Текущее содержимое cache в Redis можно посмотреть через веб-интерфейс http://localhost:8081
![screen-redis-ui.png](assets/screen-redis-ui.png)

#### Kafka UI

Посмотреть сообщения в Kafka можно через веб-интерфейс http://localhost:8383/ui/clusters/local/all-topics
![screen-kafka-ui.png](assets/screen-kafka-ui.png)

В топике [links-created](http://localhost:8383/ui/clusters/local/all-topics/links-created) содержатся информация о всех созданных коротких ссылках.

Сервис получает сообщения из [links-requested](http://localhost:8383/ui/clusters/local/all-topics/links-requested) и автоматически создает короткую ссылку при получении нового сообщения.

## Metrics

Посмотреть метрики сервиса можно в Grafana: http://localhost:3000/d/golang-metrics-dashboard/golang-metrics

![screen-grafana-dashboard.png](assets/screen-grafana-dashboard.png)

## Traces

Посмотреть трейсы сервиса можно в Jaeger: http://localhost:16686

![screen-jaeger-1.png](assets/screen-jaeger-1.png)

![screen-jaeger-2.png](assets/screen-jaeger-2.png)

## Environment variables

| Name                  | Type   | Expected | Default       | Description                              |
|-----------------------|--------|----------|---------------|------------------------------------------|
| APP_NAME              | string |          | url-shortener | service name                             |
| APP_VERSION           | string |          | 0.0.0         | service version                          |
| APP_ENV               | string |          | DEV           | service environment (DEV, PROD, etc)     |
| LOGGER_LEVEL          | string |          | error         | logging level (debug, info, warn, error) |
| LOGGER_PRETTY_CONSOLE | bool   |          | false         | logging format (text/json)               |
| SENTRY_DSN            | string |          |               | sentry DSN (disabled if empty)           |
