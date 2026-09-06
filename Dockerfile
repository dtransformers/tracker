# syntax=docker/dockerfile:1.6

############################
# 1. Build stage
############################
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies and create a non-root user
RUN apk add --no-cache git ca-certificates tzdata && \
    adduser -D -g '' appuser

# Enable reproducible builds & module caching
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Cache dependencies first (better layer caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source explicitly instead of recursively
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY pkg/ ./pkg/
COPY docs/ ./docs/
COPY tools/ ./tools/



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
    -o /out/app ./cmd/api


############################
# 2. Runtime stage
############################
FROM gcr.io/distroless/static-debian12:nonroot

# distroless/static:nonroot already includes CA certs and runs as nonroot (uid 65532)
# No need to copy /etc/passwd or manually set USER

# Copy binary
COPY --from=builder /out/app /app

ENV TMPDIR=/tmp
VOLUME ["/tmp"]

ENV LOG_DIR=/logs
VOLUME ["/logs"]

ENV APP_ENV=production

EXPOSE 8080
## TO GENERATE SWAGGER DOCS:
## go run github.com/swaggo/swag/cmd/swag init -g cmd/api/main.go      
## go run github.com/swaggo/swag/cmd/swag init -g cmd/api/main.go --parseDependency --parseInternal

ENTRYPOINT ["/app"]