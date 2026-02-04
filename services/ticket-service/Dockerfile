# ---- build stage ----
FROM golang:1.25 AS builder

WORKDIR /app

# cache deps
COPY go.mod go.sum ./
RUN go mod download

# copy source
COPY . .

# build
RUN CGO_ENABLED=0 GOOS=linux go build -o ticket-service ./cmd

# ---- run stage ----
FROM alpine:3.20

WORKDIR /app
COPY --from=builder /app/ticket-service /app/ticket-service

EXPOSE 8082
ENTRYPOINT ["/app/ticket-service"]
