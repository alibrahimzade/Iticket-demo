# ---- build stage ----
FROM golang:1.25 AS builder

WORKDIR /app

# cache deps
COPY go.mod go.sum ./
RUN go mod download

# copy source
COPY . .

# build
RUN CGO_ENABLED=0 GOOS=linux go build -o event-service ./cmd

# ---- run stage ----
FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /app/event-service /app/event-service

COPY --from=builder /app/migrations /app/migrations

EXPOSE 8081

ENTRYPOINT ["/app/event-service"]
