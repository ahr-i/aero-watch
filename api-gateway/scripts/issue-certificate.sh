#!/usr/bin/env sh
set -eu

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
ENV_FILE="${ROOT_DIR}/.env"
CERTBOT_CONF_DIR="${ROOT_DIR}/certbot/conf"

if [ ! -f "$ENV_FILE" ]; then
    echo "Missing .env file: $ENV_FILE" >&2
    echo "Create it from .env.example first." >&2
    exit 1
fi

set -a
. "$ENV_FILE"
set +a

if [ -z "${GATEWAY_SERVER_NAME:-}" ]; then
    echo "Missing required env: GATEWAY_SERVER_NAME" >&2
    exit 1
fi

if [ -z "${GATEWAY_CERTBOT_EMAIL:-}" ]; then
    echo "Missing required env: GATEWAY_CERTBOT_EMAIL" >&2
    echo "Add it to .env, for example: GATEWAY_CERTBOT_EMAIL=admin@example.com" >&2
    exit 1
fi

mkdir -p "$CERTBOT_CONF_DIR"

docker run --rm -it \
    -v "${CERTBOT_CONF_DIR}:/etc/letsencrypt" \
    certbot/certbot certonly \
    --manual \
    --preferred-challenges dns \
    --domain "$GATEWAY_SERVER_NAME" \
    --email "$GATEWAY_CERTBOT_EMAIL" \
    --agree-tos \
    --no-eff-email

echo
echo "Certificate files should now exist under:"
echo "  ${CERTBOT_CONF_DIR}/live/${GATEWAY_SERVER_NAME}/fullchain.pem"
echo "  ${CERTBOT_CONF_DIR}/live/${GATEWAY_SERVER_NAME}/privkey.pem"