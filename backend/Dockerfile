# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Set GOPROXY for faster and more reliable module downloads
ENV GOPROXY=https://goproxy.cn,direct

# Use Aliyun mirror for Alpine packages
RUN sed -i 's#http://dl-cdn.alpinelinux.org/#https://mirrors.aliyun.com/#g' /etc/apk/repositories

# Install git (required for swag)
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Install swag CLI and add to PATH
RUN go install github.com/swaggo/swag/cmd/swag@latest
ENV PATH="/root/go/bin:${PATH}"

# Copy source code
COPY . .

# Generate Swagger documentation
RUN swag init

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kaleidoscope .

# Final stage
FROM alpine:latest

# Use Aliyun mirror for Alpine packages
RUN sed -i 's#http://dl-cdn.alpinelinux.org/#https://mirrors.aliyun.com/#g' /etc/apk/repositories

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