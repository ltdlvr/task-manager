GOC=go

.PHONY: run-rest db-migrate 

run-rest:
	$(GOC) run ./cmd/rest/main.go

db-migrate:
	$(GOC) run ./cmd/db-migrate/main.go
