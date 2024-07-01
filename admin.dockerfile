FROM oven/bun:1.1.17-debian@sha256:fe1de05bfc9f648e49c41a54edd7418709e5284ca610529e11832ddb7a4870a3 AS bun

WORKDIR /app

COPY package.json bun.lockb tailwind.config.js tsconfig.json style.css ./

COPY ./internal/template/*.templ ./internal/template/

RUN bun install

RUN bunx tailwindcss -i ./style.css -o ./public/style.css


FROM golang:1.22-bullseye@sha256:067c5c7fe6d79f900c5ebe8351166356d6e3bbfcc6f807030e89b9a929252273 AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download


COPY . .

COPY --from=bun /app/public/style.css /app/public/style.css

RUN go install github.com/a-h/templ/cmd/templ@v0.2.680

RUN templ generate

RUN go build \
  -ldflags="-linkmode external -extldflags -static -X 'main.BUILDTIME=$(date --iso-8601=seconds --utc)'" \
  -o admin \
  ./cmd/admin/main.go

RUN useradd -u 1001 dxta


FROM scratch

WORKDIR /

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=build /etc/passwd /etc/passwd

COPY --from=build /app/admin /admin

USER dxta

EXPOSE 80

EXPOSE 443

CMD ["/admin"]
