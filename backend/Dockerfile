# syntax=docker/dockerfile:1.7

FROM golang:1.26-alpine3.24 AS build

WORKDIR /src

COPY go.mod ./
RUN go mod download

COPY cmd ./cmd

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/apu-api \
    ./cmd/apu-api

FROM alpine:3.24 AS runtime

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S apu \
    && adduser -S -G apu -u 10001 apu

COPY --from=build /out/apu-api /usr/local/bin/apu-api

USER apu
EXPOSE 8080

ENV APU_ENV=local \
    APU_HTTP_ADDR=:8080

HEALTHCHECK --interval=10s --timeout=3s --start-period=5s --retries=5 \
    CMD wget -qO- http://127.0.0.1:8080/healthz >/dev/null || exit 1

ENTRYPOINT ["/usr/local/bin/apu-api"]
