FROM ubuntu:22.04 AS builder

ARG GO_VERSION=1.24.0

ENV DEBIAN_FRONTEND=noninteractive
ENV PATH=/usr/local/go/bin:$PATH
ENV GOTOOLCHAIN=local

WORKDIR /src

RUN apt-get update \
	&& apt-get install -y --no-install-recommends ca-certificates wget tar \
	&& rm -rf /var/lib/apt/lists/* \
	&& wget -q "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz" -O /tmp/go.tar.gz \
	&& tar -C /usr/local -xzf /tmp/go.tar.gz \
	&& rm /tmp/go.tar.gz

COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./
COPY db ./db
COPY handler ./handler
COPY setting ./setting
COPY utils ./utils

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/drone-operation .

FROM ubuntu:22.04

ENV DEBIAN_FRONTEND=noninteractive
WORKDIR /app

RUN apt-get update \
	&& apt-get install -y --no-install-recommends ca-certificates \
	&& rm -rf /var/lib/apt/lists/* \
	&& useradd -r -u 10001 appuser

COPY --from=builder /out/drone-operation /app/drone-operation
COPY setting ./setting

EXPOSE 40002

USER appuser

CMD ["/app/drone-operation"]
