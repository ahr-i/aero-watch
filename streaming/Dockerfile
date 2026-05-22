FROM ubuntu:22.04 AS builder

ARG GO_VERSION=1.24.0
ARG TARGETARCH=amd64

ENV DEBIAN_FRONTEND=noninteractive
ENV PATH="/usr/local/go/bin:${PATH}"

RUN set -eux; \
    apt-get update; \
    apt-get install -y --no-install-recommends ca-certificates curl git; \
    case "${TARGETARCH}" in \
        amd64) GO_ARCH="amd64" ;; \
        arm64) GO_ARCH="arm64" ;; \
        *) echo "unsupported TARGETARCH=${TARGETARCH}" >&2; exit 1 ;; \
    esac; \
    curl -fsSL "https://go.dev/dl/go${GO_VERSION}.linux-${GO_ARCH}.tar.gz" -o /tmp/go.tgz; \
    tar -C /usr/local -xzf /tmp/go.tgz; \
    rm -rf /var/lib/apt/lists/* /tmp/go.tgz

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH="${TARGETARCH}" go build -o /out/streaming .

FROM ubuntu:22.04

ENV DEBIAN_FRONTEND=noninteractive

RUN set -eux; \
    apt-get update; \
    apt-get install -y --no-install-recommends ca-certificates ffmpeg tzdata; \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /out/streaming /app/streaming
COPY setting /app/setting
COPY admin /app/admin

RUN mkdir -p /app/hls /app/captures

EXPOSE 40004 41935

CMD ["/app/streaming"]
