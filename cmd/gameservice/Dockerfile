FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o gameservice ./cmd/gameservice

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/gameservice /app/gameservice
COPY --from=builder /app/configs /app/configs

# (opcionalno) certovi ako ikad treba TLS
RUN apk add --no-cache ca-certificates

EXPOSE 8081

CMD ["/app/gameservice"]
