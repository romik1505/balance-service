CURRENT_DIR = $(shell pwd)
LOCAL_BIN=$(CURRENT_DIR)/bin
VALUES=$(CURRENT_DIR)/k8s/values_local.yaml

ifndef PG_DSN
$(eval PG_DSN=$(shell cat $(VALUES) | grep -i "pg_dsn" -A1 | sed -n '2p;2q' | sed -e 's/[ \t]*value://g'))
endif

run:
	go run cmd/main.go

build:
	@go mod tidy
	CGO_ENABLED=0 go build -o bin/main cmd/main.go 

bin-deps:
	@mkdir -p bin
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.5.3

db\:up:
	$(LOCAL_BIN)/goose -dir migrations postgres "$(PG_DSN)" up

db\:down:
	$(LOCAL_BIN)/goose -dir migrations postgres "$(PG_DSN)" down

db\:create:
	$(LOCAL_BIN)/goose -dir migrations create "$(NAME)" sql

swag:
	swag init -g cmd/main.go

test:
	PG_DSN=$(PG_DSN) VALUES=$(VALUES) go test -v ./...

lint:
	$(LOCAL_BIN)/golangci-lint run ./... --timeout 60s

mocks:
	mockgen -source=./internal/service/balance/balance.go -destination=./pkg/mock/service/mock_balance/mock_service.go
	mockgen -source=./internal/service/currency/currency.go -destination=./pkg/mock/service/mock_currency/mock_currency.go
	mockgen -source=./internal/store/store.go -destination=./pkg/mock/store/mock_storage/mock_storage.go

docker-exec: bin-deps db\:up run

docker-start:
	sudo docker compose build --no-cache
	sudo docker compose up 
