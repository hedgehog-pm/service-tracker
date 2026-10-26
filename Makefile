.PHONY: fmt lint generate generate-api generate-sql

export DATABASE_URL ?= postgres://tracker:tracker@localhost:5433/tracker?sslmode=disable
export PSQL_DSN ?= $(DATABASE_URL)

lint:
	go tool golangci-lint config verify
	go tool golangci-lint run

fmt:
	go tool golangci-lint fmt

.PHONY: fix
fix:
	go tool golangci-lint config verify
	go tool golangci-lint run --fix ./...

generate: generate-sql generate-api

generate-sql:
	@test -n "$$PSQL_DSN" || (echo 'PSQL_DSN must point to a migrated PostgreSQL database'; exit 1)
	@set -eu; \
		trap 'rm -f internal/generated/bob/task/queries.sql internal/generated/bob/project/queries.sql internal/generated/bob/status/queries.sql internal/generated/bob/module/queries.sql internal/generated/bob/workflow/queries.sql internal/generated/bob/priority/queries.sql' EXIT; \
		for part in task project status module workflow priority; do \
			mkdir -p internal/generated/bob/$$part; \
			cp internal/repository/$$part/queries.sql internal/generated/bob/$$part/queries.sql; \
		done; \
		go tool bobgen-psql -c bobgen.yaml

generate-api:
	@for part in common projects statuses modules tasks transitions task-types priorities; do \
		go tool oapi-codegen -config api/server/$$part/oapi-codegen.yaml api/server/$$part/openapi.yaml || exit $$?; \
	done
