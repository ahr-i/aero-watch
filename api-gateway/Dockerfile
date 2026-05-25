FROM ubuntu:22.04

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        ca-certificates \
        gettext-base \
        nginx \
    && rm -rf /var/lib/apt/lists/* \
    && rm -f /etc/nginx/sites-enabled/default /etc/nginx/sites-available/default

COPY nginx/templates/default.conf.template /etc/nginx/templates/default.conf.template

EXPOSE 443

CMD ["sh", "-c", "envsubst '${GATEWAY_UPSTREAM} ${GATEWAY_SERVER_NAME} ${GATEWAY_SSL_CERTIFICATE} ${GATEWAY_SSL_CERTIFICATE_KEY} ${GATEWAY_CLIENT_MAX_BODY_SIZE} ${GATEWAY_CORS_ALLOW_ORIGIN} ${GATEWAY_CORS_ALLOW_METHODS} ${GATEWAY_CORS_ALLOW_HEADERS}' < /etc/nginx/templates/default.conf.template > /etc/nginx/conf.d/default.conf && nginx -t && exec nginx -g 'daemon off;'"]
