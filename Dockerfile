####################################################################
# Builder Stage
####################################################################
FROM golang:alpine AS builder

LABEL MAINTAINER="obiorapaschalugwu@gmail.com"

WORKDIR /go/src/github.com/ong-gtp/go-stockbot

# Install git for Go dependencies
RUN apk add --no-cache git

# Copy dependency files first for better Docker layer caching
COPY go.mod go.sum ./

# Download dependencies without modifying go.mod/go.sum
RUN go mod download

# Copy application source
COPY . .

# Build a static Linux binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o gostockbot


####################################################################
# Final Stage
####################################################################
FROM alpine:3.22

# Install CA certificates because Stockbot calls Stooq over HTTPS
RUN apk add --no-cache ca-certificates

# Create a dedicated non-root user
RUN addgroup -S appgroup && \
    adduser -S appuser -G appgroup

WORKDIR /app

# Copy application binary
COPY --from=builder \
    /go/src/github.com/ong-gtp/go-stockbot/gostockbot \
    /app/gostockbot

# Give the application to the non-root user
RUN chown appuser:appgroup /app/gostockbot

# Run as non-root
USER appuser

# Stockbot does not expose an HTTP API.
# It communicates through RabbitMQ.
#
# No EXPOSE is required.

ENTRYPOINT ["/app/gostockbot"]