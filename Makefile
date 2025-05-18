ifneq (,$(wildcard ./.env))
    include .env
    export
endif

.PHONY : assets generate serve

assets:
	@echo "\033[0;31mBuilding assets...\033[0m"
	rm -r static/dist || true
	@echo "\033[0;31mBuilding tailwind...\033[0m"
	npx tailwindcss -i static/tailwind.css -o static/dist/tailwind.css
	@echo "\033[0;31mBuilding htmx...\033[0m"
	cp node_modules/htmx.org/dist/htmx.min.js static/dist/htmx.min.js
	@echo "\033[0;31mBuilding htmx-sse...\033[0m"
	cp node_modules/htmx-ext-sse/sse.js static/dist/htmx-ext-sse.js
	@echo "\033[0;31mBuilding alpinejs...\033[0m"
	cp node_modules/alpinejs/dist/cdn.min.js static/dist/alpinejs.min.js
	@echo "\033[0;31mBuilding hyperscript...\033[0m"
	cp node_modules/hyperscript.org/dist/_hyperscript.min.js static/dist/_hyperscript.min.js
	@echo "\033[0;31mBuilding fonts...\033[0m"
	cp -r static/fonts static/dist/fonts

build_dungeons:
	cd dungeons && npm run build

generate:
	$(MAKE) assets
	go generate ./...

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
	