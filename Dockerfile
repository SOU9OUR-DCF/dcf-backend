FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api/main.go

FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/main .
COPY --from=builder /app/swagger.yaml .

RUN mkdir -p /app/config

COPY --from=builder /app/internal/adapters/config/config.yaml /app/config/config.yaml
COPY --from=builder /app/internal/adapters/config/config.prod.yaml /app/config/config.prod.yaml
COPY --from=builder /app/internal/adapters/config/config.dev.yaml /app/config/config.dev.yaml

RUN adduser -D -g '' appuser
RUN chown -R appuser:appuser /app
USER appuser

EXPOSE 80

CMD ["./main"]