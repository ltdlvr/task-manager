GOC=go

.PHONY: run-rest db-migrate test-api

run-rest:
	$(GOC) run ./cmd/rest/main.go

db-migrate:
	MIGRATIONS_DIR=./migrations $(GOC) run ./cmd/db-migrate/main.go

test-api:
	bash ./scripts/api-test.sh
