POSTGRES_SETUP := user=postgres password=postgres dbname=postgres host=localhost port=5432 sslmode=disable
MIGRATION_FOLDER=$(CURDIR)/pkg/db/migrations
MIGRATION_NAME=migration

.PHONY: gen-tls-cert
gen-tls-cert:
	openssl genrsa -out server.key 2048  # Сгенерировать приватный ключ (.key)
	openssl req -new -x509 -sha256 -key server.key -out server.crt -days 365 -nodes  # Сгенерировать публичный ключ (.crt), но основе приватного
	mv -f server.key server.crt configs/  # Поместить оба файла в папку /configs

.PHONY: compose-up
compose-up:
	docker compose build
	docker compose up -d postgres

.PHONY: compose-down
compose-down:
	docker compose down

.PHONY: run
run:
	go run cmd/server/main.go

.PHONY: migration-create
migration-create:
	goose -dir "$(MIGRATION_FOLDER)" create "$(MIGRATION_NAME)" sql

.PHONY: migration-up
migration-up:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP)" up

.PHONY: migration-down
migration-down:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP)" down
