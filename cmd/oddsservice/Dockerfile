# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Go module files
COPY go.mod go.sum ./
RUN go mod download

# Copy ostatka koda
COPY . .

# Build oddsservice binarke
RUN CGO_ENABLED=0 GOOS=linux go build -o oddsservice ./cmd/oddsservice

# Runtime stage
FROM alpine:3.20

WORKDIR /app

# Binarka + configi
COPY --from=builder /app/oddsservice /app/oddsservice
COPY --from=builder /app/configs /app/configs

RUN apk add --no-cache ca-certificates

EXPOSE 8083

CMD ["/app/oddsservice"]
