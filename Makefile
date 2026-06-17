GOC=go
ENV_FILE ?= .env.dev

ifneq (,$(wildcard $(ENV_FILE)))
include $(ENV_FILE)
export
endif

.PHONY: run-rest db-migrate test-api db-up

db-up:
	docker-compose up -d
run-rest:
	$(GOC) run ./cmd/rest/main.go

db-migrate:
	MIGRATIONS_DIR=./migrations $(GOC) run ./cmd/db-migrate/main.go

test-api:
	bash ./scripts/api-test.sh
