ifneq (,$(wildcard ./.env))
    include .env
    export
endif

.PHONY : generate serve

build_dungeons:
	cd dungeons && npm run build

build_www:
	cd www && npm run build

generate:
	go generate ./...
	go run cmd/gen/main.go

compose:
	docker compose -f infrastructure/docker-compose.yml --env-file=.env $(filter-out $@,$(MAKECMDGOALS))

serve:
	air -c .air.toml

db-migrate:
	@echo "\033[0;31mRunning database migrations...\033[0m"
	migration_name=$(shell date +%Y%m%d%H%M%S) && \
		atlas migrate diff ${migration_name} --dir "file://ent/migrate/migrations" --to "ent://ent/schema" --dev-url "docker://postgres/17/test?search_path=public"

database_url = "postgres://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_DATABASE}?sslmode=disable&search_path=public"

db-apply:
	@echo "\033[0;31mApplying database migrations...\033[0m"
	atlas schema apply --env local --url $(database_url)
	