.PHONY: watch-oss-api
watch-oss-api:
	@export $$(cat .env | xargs) && \
    ./bin/air -c oss-api.air.toml

.PHONY: watch-internal-api
watch-internal-api:
	@export $$(cat .env | xargs) && \
    ./bin/air -c internal-api.air.toml

.PHONY: onboarding-worker
onboarding-worker:
	@export $$(cat .env | xargs) && \
	go run cmd/onboarding-worker/main.go


.ONESHELL:
setup:
	@curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s

.PHONY: test
test:
	@go test -v ./...

.PHONY: bench
bench:
	@go test -bench=. ./...
