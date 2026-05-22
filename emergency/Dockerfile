FROM golang:1.24 AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o emergency-server .

FROM ubuntu:22.04

WORKDIR /app

RUN apt-get update \
	&& apt-get install -y --no-install-recommends ca-certificates \
	&& rm -rf /var/lib/apt/lists/*

COPY --from=builder /src/emergency-server ./emergency-server
COPY setting ./setting
COPY csv ./csv

EXPOSE 40005

CMD ["./emergency-server"]
