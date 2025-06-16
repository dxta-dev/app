FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/worker ./cmd/worker/main.go

FROM alpine:latest

RUN apk add --no-cache ca-certificates

COPY --from=builder /go/bin/worker /usr/local/bin/worker

RUN chmod +x /usr/local/bin/worker

ENTRYPOINT ["/usr/local/bin/worker"]
