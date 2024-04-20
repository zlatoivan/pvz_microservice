POSTGRES_SETUP := user=postgres password=postgres dbname=postgres host=localhost port=5433 sslmode=disable
POSTGRES_SETUP_TEST := user=postgres password=postgres dbname=test host=localhost port=5431 sslmode=disable
MIGRATION_FOLDER=$(CURDIR)/migrations
MIGRATION_NAME=pvz

.PHONY: help
help: ## display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: compose-up
compose-up: ### run docker compose
	docker compose up --build

.PHONY: compose-up-local
compose-up-local: ### run docker compose local
	docker compose up pg_db pg_db_test zookeeper kafka2 redis -d

.PHONY: compose-down
compose-down: ### down docker compose
	docker compose down


# ---------- migrations ----------

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
migration-up-test: ### migration up test
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" up

.PHONY: migration-down-test
migration-down-test: ### migration down test
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" down


# ---------- tests ----------

.PHONY: test
test: ### run tests
	go test -v -count=2 -p=3 ./...

.PHONY: test-integration
test-integration: ### run integration tests
	CONFIG_PATH=$(CURDIR)/config/config_test.yaml go test ./... -tags=integration -v

.PHONY: run-test
run-test: ### run with tests
	CONFIG_PATH=config/config_test.yaml go run cmd/server/main.go


# ---------- local ----------

.PHONY: run
run: ### run local
	CONFIG_PATH=config/config.yaml go run cmd/server/main.go

..PHONY: build
build: ### build local
	go build -o main cmd/server/main.go


# ---------- rare ----------

.PHONY: linter
linter: ### check by golangci linter
	golangci-lint run

.PHONY: docker-rm-volume
docker-rm-volume: ### remove docker volume
	docker volume rm pg-data

.PHONY: gen-ssl-cert
gen-ssl-cert: ### generate fresh ssl certificate
	openssl genrsa -out server.key 2048  # Сгенерировать приватный ключ (.key)
	openssl req -new -x509 -sha256 -key server.key -out server.crt -days 365 -nodes  # Сгенерировать публичный ключ (.crt), но основе приватного
	mv -f server.key server.crt internal/server/certs/  # Поместить оба файла в папку /certificate



# ---------- tests with grc ----------

.PHONY: test-grc
test-grc: ### run tests grc
	grc go test -v -count=2 -p=3 ./...

.PHONY: test-integration-grc
test-integration-grc: ### run integration tests grc
	CONFIG_PATH=$(CURDIR)/config/config_test.yaml grc go test ./... -tags=integration -v