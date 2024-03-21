POSTGRES_SETUP := user=postgres password=postgres dbname=postgres host=localhost port=5432 sslmode=disable
MIGRATION_FOLDER=$(CURDIR)/pkg/db/migrations
MIGRATION_NAME=migration

.PHONY: migration-create
migration-create:
	goose -dir "$(MIGRATION_FOLDER)" create "$(MIGRATION_NAME)" sql

.PHONY: migration-up
migration-up:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP)" up

.PHONY: migration-down
migration-down:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP)" down
