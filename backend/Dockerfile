# =============================================================================
# Multi-stage Dockerfile for aLMS API
# =============================================================================

# -----------------------------------------------------------------------------
# Stage 1: Build
# -----------------------------------------------------------------------------
FROM golang:1.25-alpine AS builder

# Install git and ca-certificates (needed for fetching dependencies)
RUN apk add --no-cache git ca-certificates tzdata

# Create non-root user for runtime
RUN adduser -D -g '' appuser

WORKDIR /build

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the binary with optimizations
ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.Version=${VERSION}" \
    -o /build/api \
    ./cmd/api

# Build migrate binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /build/migrate \
    ./cmd/migrate

# -----------------------------------------------------------------------------
# Stage 2: Runtime (minimal)
# -----------------------------------------------------------------------------
FROM scratch

# Import from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/passwd /etc/passwd

# Copy binaries
COPY --from=builder /build/api /api
COPY --from=builder /build/migrate /migrate

# Copy migrations folder (needed for migrate container)
COPY --from=builder /build/migrations /migrations

# Run as non-root
USER appuser

# Configurable port (default 8080)
ENV PORT=8080
EXPOSE ${PORT}

# Run the API
ENTRYPOINT ["/api"]
