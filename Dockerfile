FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the Go app
RUN go build -o bsf-engine main.go

# Final image
FROM alpine:latest

WORKDIR /app

# Copy built binary and assets
COPY --from=builder /app/bsf-engine .
COPY assets ./assets
COPY .env.template ./

# Expose port (default 9090, can be overridden)
EXPOSE 9090

# Set environment defaults (can be overridden at runtime)
ENV PORT=9090
ENV REDIS_HOST=localhost
ENV REDIS_PORT=6379

# Run the app
CMD ["./bsf-engine"]