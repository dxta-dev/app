.PHONY: watch-admin
watch-admin: CONFIG := admin.air.toml
watch-admin: _watch

.PHONY: watch-web
watch-web: CONFIG := web.air.toml
watch-web: _watch

_watch:
	@export $$(cat .env | xargs) && \
    ./bin/air -c $(CONFIG) & \
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
	@go install github.com/a-h/templ/cmd/templ@v0.2.747 && cp $(shell go env GOPATH)/bin/templ ./bin
	@bun i

.PHONY: test
test:
	@go test -v ./...

.PHONY: bench
bench:
	@go test -bench=. ./...
