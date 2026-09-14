# === STAGE 1: COMPILATION ENGINE ===
FROM golang:1.24-alpine AS builder

# Installing gcc, musl-dev, and sqlite-dev is critical because go-sqlite3 uses CGO under the hood!
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app
RUN mkdir -p /app/data

# Pre-fetch and cache Go package dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire workspace source tree into the container build frame
COPY . .

# Build a fully optimized, statically linked CGO production binary
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-extldflags=-static" -o forum-exec ./cmd/web/main.go

# === STAGE 2: MINIMAL RUNTIME CONTAINER ===
FROM alpine:3.19

RUN apk add --no-cache ca-certificates

WORKDIR /app

RUN mkdir -p /app/data

# Copy the compiled production executable from the builder stage
COPY --from=builder /app/forum-exec .

# Copy structural layout assets and the database migrations schema
COPY --from=builder /app/ui ./ui
COPY --from=builder /app/internal/database/schema.sql ./internal/database/schema.sql

# Expose the network communication port
EXPOSE 8080

# Configure environment variables to locate the local embedded SQLite file directory path safely
ENV DB_PATH=/app/forum.db
ENV SCHEMA_PATH=/app/internal/database/schema.sql

# Execute the application
CMD ["./forum-exec"]
