# Backend Dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app

# Install git and build dependencies with pinned versions
RUN apk add --no-cache git
# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy swagger docs if they exist
COPY docs ./docs

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main main.go

# Final stage
FROM alpine:3.19
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .
# Copy swagger docs if they exist
COPY --from=builder /app/docs ./docs

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./main"]

