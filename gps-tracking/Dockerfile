FROM ubuntu:22.04 AS builder

ARG GO_VERSION=1.24.0
ARG TARGETARCH=amd64

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates wget \
    && rm -rf /var/lib/apt/lists/*

RUN wget -q "https://go.dev/dl/go${GO_VERSION}.linux-${TARGETARCH}.tar.gz" -O /tmp/go.tar.gz \
    && tar -C /usr/local -xzf /tmp/go.tar.gz \
    && rm /tmp/go.tar.gz

ENV PATH="/usr/local/go/bin:${PATH}"
ENV CGO_ENABLED=0
ENV GOWORK=off

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -trimpath -ldflags="-s -w" -o /out/gps-tracking .

FROM ubuntu:22.04

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /out/gps-tracking ./gps-tracking
COPY setting ./setting

EXPOSE 40003

CMD ["./gps-tracking"]
