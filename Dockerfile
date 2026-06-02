# syntax=docker/dockerfile:1.6

############################
# 1. Build stage
############################
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies (if needed)
RUN apk add --no-cache git ca-certificates

# Enable reproducible builds & module caching
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Cache dependencies first (better layer caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build-time metadata (inject at build time)
ARG VERSION=dev
ARG COMMIT=none
ARG BUILD_TIME=unknown

# Build binary (static)
RUN go build -trimpath -ldflags "\
    -s -w \
    -X main.version=${VERSION} \
    -X main.commit=${COMMIT} \
    -X main.buildTime=${BUILD_TIME}" \
    -o /out/app ./cmd/app




############################
# 2. Runtime stage
############################
FROM scratch

# CA certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy binary
COPY --from=builder /out/app /app


ENV TMPDIR=/tmp
VOLUME ["/tmp"]

ENV LOG_DIR=/logs
VOLUME ["/logs"]

ENV APP_ENV=production

EXPOSE 8080

ENTRYPOINT ["/app"]