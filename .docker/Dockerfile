# ===============================
# Stage 1: Builder
# ===============================
FROM golang:1.25-alpine AS builder

# Install git and build tools
RUN apk add --no-cache git build-base

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum, then download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code
COPY . .

# Ensure the tmp folder has proper permissions (some builds may need it)
RUN mkdir -p /app/tmp && chmod -R 777 /app/tmp

# Build the Go binary for Linux (disable CGO for a static binary)
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -v -o /app/main .

# ===============================
# Stage 2: Production Runtime
# ===============================
FROM alpine:latest

# Install SSL root certificates and timezone data
RUN apk --no-cache add ca-certificates tzdata

# Set timezone environment variable
ENV TZ=Asia/Jakarta

# Copy timezone data to a known location (optional)
RUN cp /usr/share/zoneinfo/Asia/Jakarta /etc/localtime || true

# Set working directory
WORKDIR /app

# Copy the pre-built binary from the builder stage
COPY --from=builder /app/main .

# Create tmp folder in runtime image (if the app requires it)
RUN mkdir -p /app/tmp && chmod -R 777 /app/tmp

# Expose the application port
EXPOSE 8000

# Command to start the application
CMD ["./main"]
