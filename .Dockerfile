# Build stage
FROM --platform=linux/amd64 golang:1.23-alpine AS builder

# Add only necessary build tools and add security
RUN apk add --no-cache gcc musl-dev && \
    adduser -D -u 10001 appuser

WORKDIR /app

# Copy only files needed for dependency resolution first (better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build with optimized flags
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -a \
    -ldflags="-s -w -extldflags '-static'" \
    -o binance-mark-prices-sync \
    ./cmd

# Final stage - using distroless for better security
FROM gcr.io/distroless/static-debian12:nonroot

# Copy binary from builder
WORKDIR /app
COPY --from=builder /app/alerts-worker /app/

# Use nonroot user
USER 10001:10001

# Expose any necessary ports (uncomment if needed)
# EXPOSE 8080

# Set the binary as the entrypoint
ENTRYPOINT ["/app/alerts-worker"]