FROM ubuntu:22.04 AS builder

ARG GO_VERSION=1.24.0
ARG DEBIAN_FRONTEND=noninteractive

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates wget tar git \
    && rm -rf /var/lib/apt/lists/*

RUN wget -q "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz" -O /tmp/go.tar.gz \
    && tar -C /usr/local -xzf /tmp/go.tar.gz \
    && rm /tmp/go.tar.gz

ENV PATH="/usr/local/go/bin:${PATH}"
ENV CGO_ENABLED=0

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -trimpath -ldflags="-s -w" -o /out/auth .

FROM ubuntu:22.04

ARG DEBIAN_FRONTEND=noninteractive
ENV TZ=UTC

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates tzdata \
    && rm -rf /var/lib/apt/lists/* \
    && useradd -r -u 10001 -g nogroup auth

WORKDIR /app

COPY --from=builder /out/auth ./auth
COPY --from=builder /src/setting ./setting

USER auth

EXPOSE 40001

ENTRYPOINT ["./auth"]
