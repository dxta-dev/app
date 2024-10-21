FROM oven/bun:1.1.31-debian@sha256:e0139b1f0f2e96c2bc39737e4f74837e84b6248b8939676c2490416d5b6f0b1d AS bun

WORKDIR /app

COPY package.json bun.lockb tailwind.config.js tsconfig.json style.css ./

COPY ./internal/template/*.templ ./internal/template/

RUN bun install

RUN bunx tailwindcss -i ./style.css -o ./public/style.css


FROM golang:1.23-bullseye@sha256:ecef8303ced05b7cd1addf3c8ea98974f9231d4c5a0c230d23b37bb623714a23 AS build

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
