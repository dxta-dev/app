AIR_CONFIG ?= admin.air.toml

.PHONY: watch
watch:
	@export $$(cat .env | xargs) && \
    ./bin/air -c $(AIR_CONFIG) & \
	$(MAKE) tailwind-watch

.PHONY: templ
templ:
	@./bin/templ generate ./internal/templates/*.templ

.PHONY: tailwind-watch
tailwind-watch:
	@bunx tailwindcss -i ./style.css -o ./public/style.css --watch

.PHONY: tailwind-build
tailwind-build:
	@./node_modules/.bin/tailwindcss -i ./style.css -o ./public/style.css

.ONESHELL:
setup:
	@curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s
	@go install github.com/a-h/templ/cmd/templ@latest && cp $(shell go env GOPATH)/bin/templ ./bin
	@bun i

.PHONY: test
test:
	@go test -v ./...

.PHONY: bench
bench:
	@go test -bench=. ./...
