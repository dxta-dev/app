FROM golang:1.24-bullseye@sha256:cf29cafe674ad5e637311148fb7933f67c4e8cafa79ce066aad0e7aa708fc287 AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build \
  -ldflags="-linkmode external -extldflags -static -X 'main.BUILDTIME=$(date --iso-8601=seconds --utc)'" \
  -o api \
  ./cmd/api/main.go

RUN useradd -u 1001 dxta

FROM scratch

WORKDIR /

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=build /etc/passwd /etc/passwd

COPY --from=build /app/api /api

USER dxta

EXPOSE 80

EXPOSE 443

CMD ["/api"]
