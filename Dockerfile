# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod ./
COPY go.sum* ./

# Download dependencies (only if go.sum exists and has content)
RUN if [ -f go.sum ] && [ -s go.sum ]; then go mod download; fi

# Copy source code
COPY . .

# Build arguments for version info (optional overrides)
ARG BUILD_TIME
ARG GIT_COMMIT

# Build the application
# Note: Version is defined in main.go, only BUILD_TIME and GIT_COMMIT are injected
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT}" \
    -o /app/demo-app .

# Runtime stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/demo-app /app/demo-app

# Create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

EXPOSE 8000

ENV PORT=8000
ENV OTEL_EXPORTER_OTLP_ENDPOINT=jaeger:4318
ENV OTEL_SERVICE_NAME=demo-app
ENV APP_ENV=production

CMD ["/app/demo-app"]
