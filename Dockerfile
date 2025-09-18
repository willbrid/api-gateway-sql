FROM golang:1.24 AS builder

RUN apt-get update && apt-get install -y gcc sqlite3 libsqlite3-dev

WORKDIR /build

COPY . .

RUN go mod download

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o api-gateway-sql ./cmd

FROM debian:bookworm-slim

COPY --from=builder /build/api-gateway-sql /usr/local/bin
COPY --from=builder /build/entrypoint.sh /usr/local/bin
RUN chmod +x /usr/local/bin/entrypoint.sh

RUN mkdir /data
RUN mkdir /etc/api-gateway-sql
RUN mkdir /etc/api-gateway-sql/tls

RUN groupadd -r nobody && usermod -g nobody nobody
RUN chown -R nobody:nobody /data

COPY --from=builder --chown=nobody /build/fixtures/config.default.yaml /etc/api-gateway-sql/config.yaml
COPY --from=builder --chown=nobody /build/fixtures/tls/server.crt /etc/api-gateway-sql/tls/server.crt
COPY --from=builder --chown=nobody /build/fixtures/tls/server.key /etc/api-gateway-sql/tls/server.key

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates

ENV API_GATEWAY_SQL_PORT=5297
ENV API_GATEWAY_SQL_ENABLE_HTTPS="true"

USER nobody
EXPOSE $API_GATEWAY_SQL_PORT

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]