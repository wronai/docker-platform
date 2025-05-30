# Build stage
FROM golang:1.21-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Initialize the Go module if not already done
RUN go mod init media-vault 2>/dev/null || true

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Tidy up the module dependencies
RUN go mod tidy

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o media-vault-api ./cmd

# Final stage
FROM alpine:3.18

# Install required packages
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/media-vault-api .

# Copy migrations
COPY --from=builder /app/migrations ./migrations

# Expose the application port
EXPOSE 8080

# Command to run the application
CMD ["./media-vault-api"]