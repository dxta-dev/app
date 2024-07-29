FROM oven/bun:1.1.21-debian@sha256:c59bbb19520a5d9417c0848496758ce0d20a6054500f0916b0939540df253839 AS bun

WORKDIR /app

COPY package.json bun.lockb tailwind.config.js tsconfig.json style.css ./

COPY ./internal/template/*.templ ./internal/template/

RUN bun install

RUN bunx tailwindcss -i ./style.css -o ./public/style.css



FROM golang:1.22-bullseye@sha256:11ff2c5de8044d8d749cc684f6f3cdfbc9b2755960a62d4948b3483dc090e7aa AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download


COPY . .

COPY --from=bun /app/public/style.css /app/public/style.css

RUN go install github.com/a-h/templ/cmd/templ@v0.2.680

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
