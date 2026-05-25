# syntax=docker/dockerfile:1

FROM golang:1.24 AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/ai-analysis .

FROM ubuntu:22.04

WORKDIR /app

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates tzdata \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /out/ai-analysis ./ai-analysis
COPY setting ./setting

EXPOSE 40006

CMD ["./ai-analysis"]
