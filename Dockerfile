# --- Stage 1: Builder ---
FROM golang:1.24.4-alpine AS builder

WORKDIR /app

# Install git (needed by some Go modules)
RUN apk add --no-cache git

# Copy go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o main .

# --- Stage 2: Final Image ---
FROM alpine:latest

WORKDIR /app

# Install minimal dependencies for running Go binary
RUN apk add --no-cache ca-certificates

# Copy the compiled binary
COPY --from=builder /app/main .

# Copy static assets and templates
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static

# Copy the .env file
COPY --from=builder /app/.env .env

COPY --from=builder /app/database ./database

# Expose app port
EXPOSE 8081

# Run the app
CMD ["./main"]
