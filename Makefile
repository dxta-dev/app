.PHONY: tailwind-watch
tailwind-watch:
	@./node_modules/.bin/tailwindcss -i ./style.css -o ./public/style.css --watch

.PHONY: tailwind-build
tailwind-build:
	@./node_modules/.bin/tailwindcss -i ./style.css -o ./public/style.css

