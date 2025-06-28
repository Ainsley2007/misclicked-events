# Use the official Golang image as build environment
FROM golang:1.23 as builder

# Set working directory
WORKDIR /app

# Copy go mod files and download dependencies
COPY go.mod ./
RUN go mod download

# Copy the rest of the source
COPY . .

# Build the Go app
RUN go build -o bot ./cmd/misclickedevents

# Final minimal image
FROM debian:bookworm-slim

# Install CA certificates for HTTPS to work
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /app

# Copy built binary only
COPY --from=builder /app/bot /app/bot

# Entrypoint
CMD ["./bot"]