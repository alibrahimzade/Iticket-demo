# ---- build stage ----
FROM golang:1.25 AS builder

WORKDIR /app

# cache deps
COPY go.mod go.sum ./
RUN go mod download

# copy source
COPY . .

# build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o ticket-service ./cmd

# ---- run stage ----
FROM alpine:3.20

# 👉 Needed for HTTPS, DB TLS, future integrations
RUN apk add --no-cache ca-certificates

WORKDIR /app

# copy binary
COPY --from=builder /app/ticket-service /app/ticket-service

# 👉 COPY MIGRATIONS (IMPORTANT)
COPY --from=builder /app/migrations /app/migrations

EXPOSE 8082

ENTRYPOINT ["/app/ticket-service"]
