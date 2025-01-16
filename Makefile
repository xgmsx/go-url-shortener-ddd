export

## Installation

.PHONY: install
install: ## Install development utils
	@echo "* Running migrate-install..."
	$(MAKE) migrate-install
	@echo "* Running lint-install..."
	$(MAKE) openapi-generate
	@echo "* Running generate-install..."
	$(MAKE) generate-install

## Generation

.PHONY: generate-install
generate-install:
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install github.com/swaggo/swag/cmd/swag@latest

.PHONY: proto-generate
proto-generate:
	protoc --go_out=proto/gen --go-grpc_out=proto/gen --proto_path=proto proto/shortener_v1.proto

.PHONY: openapi-generate
openapi-generate:
	swag init --generalInfo ./pkg/http/server.go --parseInternal

.PHONY: generate
generate: ## Generate artifacts
	@echo "* Running proto-generate..."
	$(MAKE) proto-generate
	@echo "* Running openapi-generate..."
	$(MAKE) openapi-generate
	@echo "* Running go generate..."
	go generate ./...
	$(MAKE) fmt


## Migrations

DB_MIGRATE_URL = postgres://login:pass@localhost:5432/app-db?sslmode=disable
MIGRATE_PATH = ./migrations

.PHONY: migrate-install
migrate-install:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.1

.PHONY: migrate-create
migrate-create:  # usage: make migrate-create name=init
	migrate create -ext sql -dir "$(MIGRATE_PATH)" $(name)

.PHONY: migrate-up
migrate-up:
	migrate -database "$(DB_MIGRATE_URL)" -path "$(MIGRATE_PATH)" up

.PHONY: migrate-down
migrate-down:
	migrate -database "$(DB_MIGRATE_URL)" -path "$(MIGRATE_PATH)" down -all

## docker compose

.PHONY: up
up: ## 🐳🔼 Start docker containers
	docker compose up -d --build

.PHONY: down
down: ## 🐳🔽 Stop docker containers
	docker compose down

## Tests, linting, generation

.PHONY: lint-install
lint-install:
	go install mvdan.cc/gofumpt@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: lint
lint: ## 🚨 Run lint checks
	golangci-lint run --fix ./...

.PHONY: fmt
fmt: ## 🎨 Fix code format issues
	gofumpt -w .
	goimports -w -local github.com/xgmsx/go-url-shortener .

## Building, running, escape analysis

.PHONY: build
build: ## 📦 Build the program
	go build -o app.bin ./cmd/app

.PHONY: run
run: ## 🚶 Run the program
	go run ./cmd/app

.PHONY: rune
rune: ## 🔎 Run the program with escape analysis
	go run -gcflags='-m=3' ./cmd/app

## Tests, tests, tests...

.PHONY: cov
cov: ## ☔ Generate a coverage report
	go test -cover -coverprofile=coverage.txt ./...
	go tool cover -html=coverage.txt -o coverage.html
	go tool cover -func=coverage.txt

.PHONY: test
test: ## 🚦 Execute unittests
	go test ./...

.PHONY: test-race
test-race: ## 🚦🏁 Execute tests with the data race detector
	go test -race -short ./...

.PHONY: test-msan
test-msan: ## 🚦🧼 Execute tests with the memory sanitizer
	go test -msan -short ./...

.PHONY: test-bench
test-bench: ## 🚦 Execute benchmark tests
	go test -bench=. ./...

## Pprof

.PHONY: pprof-allocs
pprof-allocs: ## 📈 Show pprof memory allocation report
	go tool pprof -http=:8001 http://localhost:8000/debug/pprof/allocs?debug=1

.PHONY: pprof-heap
pprof-heap: ## 📈 Show pprof memory heap report
	go tool pprof -http=:8002 http://localhost:8000/debug/pprof/heap?debug=1

.PHONY: pprof-goroutine
pprof-goroutine: ## 📈 Show pprof memory heap report
	go tool pprof -http=:8003 http://localhost:8000/debug/pprof/goroutine?debug=1
