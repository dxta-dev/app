FROM oven/bun:1.0.17-debian AS bun

WORKDIR /app

COPY package.json bun.lockb tailwind.config.js tsconfig.json style.css ./

COPY ./internals/templates/*.templ ./internals/templates/

RUN bun install

RUN bunx tailwindcss -i ./style.css -o ./public/style.css



FROM golang:1.21.5-bullseye AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
  --mount=type=cache,target=/root/.cache/go-build \
  go mod download

RUN useradd -u 1001 dxta

COPY . .

RUN go install github.com/a-h/templ/cmd/templ@latest

RUN templ generate

RUN go build \
  -ldflags="-linkmode external -extldflags -static" \
  -tags netgo \
  -o web \
  ./cmd/web/main.go



FROM scratch

WORKDIR /

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=bun /app/public/style.css /public/style.css

COPY --from=build /app/public/fonts/* /public/fonts/

COPY --from=build /etc/passwd /etc/passwd

COPY --from=build /app/web /web

USER dxta

EXPOSE 3000

CMD ["/web"]
