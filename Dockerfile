FROM oven/bun:1.0.36-debian@sha256:083d9cb26552327af6228fde549c535117ce1be7e38e78e3b9d2a904cb307309 AS bun

WORKDIR /app

COPY package.json bun.lockb tailwind.config.js tsconfig.json style.css ./

COPY ./internal/template/*.templ ./internal/template/

RUN bun install

RUN bunx tailwindcss -i ./style.css -o ./public/style.css



FROM golang:1.22-bullseye@sha256:af93761952000beb347da16073f408508a98655ca00223154b9518f46c77bd07 AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download


COPY . .

COPY --from=bun /app/public/style.css /app/public/style.css

RUN go install github.com/a-h/templ/cmd/templ@v0.2.648

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
