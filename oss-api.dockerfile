FROM golang:1.24-bullseye@sha256:f0fe88a509ede4f792cbd42056e939c210a1b2be282cfe89c57a654ef8707cd2 AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build \
  -ldflags="-linkmode external -extldflags -static" \
  -o ./tmp/oss-api \
  ./cmd/oss-api/main.go

RUN useradd -u 1001 dxta

FROM scratch

WORKDIR /

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=build /etc/passwd /etc/passwd

COPY --from=build /app/tmp/oss-api /oss-api

USER dxta

EXPOSE 80

EXPOSE 443

CMD ["/oss-api"]
