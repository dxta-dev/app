.PHONY: watch-oss-api
watch-oss-api:
	@export $$(cat .env | xargs) && \
    ./bin/air -c oss-api.air.toml

watch-internal-api:
	@export $$(cat .env | xargs) && \
    ./bin/air -c internal-api.air.toml

.ONESHELL:
setup:
	@curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s

.PHONY: test
test:
	@go test -v ./...

.PHONY: bench
bench:
	@go test -bench=. ./...
