# Dockerfile
# Multi-stage build for LazyDnD

# Stage 1: Builder
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o lazydnd .

# Stage 2: Runtime
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -S lazydnd && adduser -S lazydnd -G lazydnd

WORKDIR /home/lazydnd

# Copy binary from builder
COPY --from=builder /app/lazydnd .

# Copy assets
COPY --from=builder /app/assets ./assets

# Change ownership
RUN chown -R lazydnd:lazydnd /home/lazydnd

# Switch to non-root user
USER lazydnd

# Set environment variables
ENV TERM=xterm-256color

# Run the application
ENTRYPOINT ["./lazydnd"]

