.PHONY: watch-api
watch-api:
	@export $$(cat .env | xargs) && \
    ./bin/air -c api.air.toml

watch-other-api:
	@export $$(cat .env | xargs) && \
    ./bin/air -c other-api.air.toml

.ONESHELL:
setup:
	@curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s

.PHONY: test
test:
	@go test -v ./...

.PHONY: bench
bench:
	@go test -bench=. ./...
