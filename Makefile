.PHONY: air-watch
watch:
	@./bin/air & $(MAKE) tailwind-watch

.PHONY: templ
templ:
	@./bin/templ generate ./internals/templates/*.templ

.PHONY: tailwind-watch
tailwind-watch:
	@./node_modules/.bin/tailwindcss -i ./style.css -o ./public/style.css --watch

.PHONY: tailwind-build
tailwind-build:
	@./node_modules/.bin/tailwindcss -i ./style.css -o ./public/style.css

.ONESHELL:
setup:
	@curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s
	@go install github.com/a-h/templ/cmd/templ@latest && cp $(shell go env GOPATH)/bin/templ ./bin
	@bun i
