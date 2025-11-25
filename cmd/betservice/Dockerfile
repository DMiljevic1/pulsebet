FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o betservice ./cmd/betservice

# 🔹 Runtime stage
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/betservice /app/betservice
COPY --from=builder /app/configs /app/configs

RUN apk add --no-cache ca-certificates

EXPOSE 8082

CMD ["/app/betservice"]
