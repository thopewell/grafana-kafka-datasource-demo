# Build stage
FROM golang:1.20-alpine AS builder

WORKDIR /app

# Copy all producer files
COPY . .

# Download dependencies and generate go.sum
RUN go mod tidy

# Build the binary
RUN go build -o /producer producer.go

# Runtime stage
FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=builder /producer /producer
ENTRYPOINT ["/producer"]
