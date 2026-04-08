# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Set GOPROXY for faster and more reliable module downloads
ENV GOPROXY=https://goproxy.cn,direct

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kaleidoscope .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/kaleidoscope .

# Create logs directory
RUN mkdir -p logs

# Expose port (production uses 9000)
EXPOSE 9000

# Run the binary
CMD ["./kaleidoscope", "server", "--port", "9000"]