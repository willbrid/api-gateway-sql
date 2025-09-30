#!/bin/sh

exec /usr/local/bin/api-gateway-sql \
  --config-file $API_GATEWAY_SQL_CONFIG_FILE \
  --port $API_GATEWAY_SQL_PORT \
  --enable-https $API_GATEWAY_SQL_ENABLE_HTTPS \
  --cert-file $API_GATEWAY_SQL_CERT_FILE \
  --key-file $API_GATEWAY_SQL_KEY_FILE