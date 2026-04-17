# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

RUN --mount=type=cache,target=/root/.cache/go/build \
    go install github.com/swaggo/swag/cmd/swag@latest
ENV PATH="/root/go/bin:${PATH}"

COPY . .

RUN swag init

RUN --mount=type=cache,target=/root/.cache/go/build \
    CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags "-s -w" \
    -o kaleidoscope .

# Final stage
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /app/kaleidoscope .

RUN mkdir -p logs

EXPOSE 9000

CMD ["./kaleidoscope", "server", "--port", "9000"]