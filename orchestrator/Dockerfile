FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o orchestrator .

FROM ubuntu:22.04

WORKDIR /app

COPY --from=builder /app/orchestrator ./orchestrator
COPY --from=builder /app/setting ./setting

EXPOSE 40008

CMD ["./orchestrator"]
