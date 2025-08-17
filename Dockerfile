# Multi-stage build to minimize final image size
FROM golang:1.25.0-alpine AS builder

# Build arguments for versioning
ARG VERSION=dev
ARG COMMIT=none
ARG DATE=unknown

# Install required packages for building
RUN apk add --no-cache git ca-certificates

# Create working directory
WORKDIR /app

# Copy go.mod for dependency caching
COPY go.mod ./
COPY go.sum ./

# Download dependencies (cached if go.mod unchanged)
RUN go mod download && go mod verify

# Copy only necessary files (better caching)
COPY *.go ./
COPY internal/ ./internal/

# Build application with versioning and optimization
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags "-s -w -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE} -X main.builtBy=docker" \
    -trimpath \
    -o goroutines-tester .

# Final image
FROM alpine:latest

# Install ca-certificates for HTTPS and curl for healthcheck
RUN apk --no-cache add ca-certificates curl

# Create user for security
RUN adduser -D -s /bin/sh appuser

WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/goroutines-tester .

# Change ownership to appuser
RUN chown appuser:appuser /app/goroutines-tester

# Switch to non-privileged user
USER appuser

EXPOSE 8888

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8888/metrics || exit 1


CMD ["./goroutines-tester"]
