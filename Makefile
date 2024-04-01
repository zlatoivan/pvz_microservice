POSTGRES_SETUP := user=postgres password=postgres dbname=postgres host=localhost port=5432 sslmode=disable
MIGRATION_FOLDER=$(CURDIR)/migrations
MIGRATION_NAME=pvz

.PHONY: compose-up
compose-up:
	#docker compose up --build
	docker compose up -d pg_db
	docker compose build

.PHONY: compose-down
compose-down:
	docker compose down

.PHONY: run
run:
	CONFIG_PATH=config/config.yaml go run cmd/server/main.go

.PHONY: migration-create
migration-create:
	goose -dir "$(MIGRATION_FOLDER)" create "$(MIGRATION_NAME)" sql

.PHONY: migration-up
migration-up:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP)" up

.PHONY: migration-down
migration-down:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP)" down

.PHONY: gen-ssl-cert
gen-ssl-cert:
	openssl genrsa -out server.key 2048  # Сгенерировать приватный ключ (.key)
	openssl req -new -x509 -sha256 -key server.key -out server.crt -days 365 -nodes  # Сгенерировать публичный ключ (.crt), но основе приватного
	mv -f server.key server.crt internal/server/certs/  # Поместить оба файла в папку /certs

..PHONY: linter
linter:
	golangci-lint run

..PHONY: build
build:
	go build cmd/server/main.go