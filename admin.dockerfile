FROM oven/bun:1.1.8-debian@sha256:e1f5bb763d299b68ed5a739e6177079efed3a5dd227507420e0a6c6a5cfa99d0 AS bun

WORKDIR /app

COPY package.json bun.lockb tailwind.config.js tsconfig.json style.css ./

COPY ./internal/template/*.templ ./internal/template/

RUN bun install

RUN bunx tailwindcss -i ./style.css -o ./public/style.css


FROM golang:1.22-bullseye@sha256:78b171fe51f25b8c3197710f281dacae94759a254e9a486576005f9dadba9e7d AS build

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
