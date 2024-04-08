POSTGRES_SETUP := user=postgres password=postgres dbname=postgres host=localhost port=5433 sslmode=disable
POSTGRES_SETUP_TEST := user=postgres password=postgres dbname=test host=pg_db_test port=5431 sslmode=disable
MIGRATION_FOLDER=$(CURDIR)/migrations
MIGRATION_NAME=pvz

.PHONY: help
help: ## display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: compose-up
compose-up: ### run docker compose
	#docker compose up --build
	docker compose up -d pg_db
	docker compose build

.PHONY: compose-down
compose-down: ### down docker compose
	docker compose down

.PHONY: docker-rm-volume
docker-rm-volume: ### remove docker volume
	docker volume rm pg-data

.PHONY: migration-create
migration-create: ### create new migration
	goose -dir "$(MIGRATION_FOLDER)" create "$(MIGRATION_NAME)" sql

.PHONY: migration-up
migration-up: ### migration up
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP)" up

.PHONY: migration-down
migration-down: ### migration down
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP)" down

.PHONY: migration-up-test
migration-up-test: ### migration up
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" up

.PHONY: gen-ssl-cert
gen-ssl-cert: ### generate fresh ssl certificate
	openssl genrsa -out server.key 2048  # Сгенерировать приватный ключ (.key)
	openssl req -new -x509 -sha256 -key server.key -out server.crt -days 365 -nodes  # Сгенерировать публичный ключ (.crt), но основе приватного
	mv -f server.key server.crt internal/server/certs/  # Поместить оба файла в папку /certificates

.PHONY: linter
linter: ### check by golangci linter
	golangci-lint run

#.PHONY: test
test: ### run tests
	go test -v ./...



# ---------- no need ----------

.PHONY: run
run:
	CONFIG_PATH=config/config.yaml go run cmd/server/main.go

..PHONY: build
build:
	go build -o main cmd/server/main.go

