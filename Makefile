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
	cp node_modules/alpinejs/dist/cdn.min.js static/dist/cdn.min.js
	@echo "\033[0;31mBuilding hyperscript...\033[0m"
	cp node_modules/hyperscript.org/dist/_hyperscript.min.js static/dist/_hyperscript.min.js
	@echo "\033[0;31mBuilding fonts...\033[0m"
	cp -r static/fonts static/dist/fonts

generate:
	$(MAKE) assets
	go generate

serve:
	air -c .air.toml