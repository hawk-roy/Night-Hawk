FROM golang:1.26.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/go-order-service ./cmd/server

FROM alpine:3.20

WORKDIR /app

RUN adduser -D -g '' appuser

COPY --from=builder /app/go-order-service /app/go-order-service

EXPOSE 9000

USER appuser

CMD ["./go-order-service"]