FROM oven/bun:1.0.26-debian AS bun

WORKDIR /app

COPY package.json bun.lockb tailwind.config.js tsconfig.json style.css ./

COPY ./internal/templates/*.templ ./internal/templates/

RUN bun install

RUN bunx tailwindcss -i ./style.css -o ./public/style.css



FROM golang:1.22-bullseye AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download


COPY . .

COPY --from=bun /app/public/style.css /app/public/style.css

RUN go install github.com/a-h/templ/cmd/templ@v0.2.543

RUN templ generate

RUN go build \
  -ldflags="-linkmode external -extldflags -static -X 'main.BUILDTIME=$(date --iso-8601=seconds --utc)'" \
  -o web \
  ./cmd/web/main.go

RUN useradd -u 1001 dxta


FROM scratch

WORKDIR /

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=build /etc/passwd /etc/passwd

COPY --from=build /app/web /web

USER dxta

EXPOSE 80

EXPOSE 443

CMD ["/web"]
