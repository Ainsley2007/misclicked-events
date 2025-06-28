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

# Add this to install root CAs
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && \
    rm -rf /var/lib/apt/lists/*


# Set working directory
WORKDIR /app

# Copy binary and assets
COPY --from=builder /app/bot /app/bot

# Set env and run
CMD ["./bot"]
