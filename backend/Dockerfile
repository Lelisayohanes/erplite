# Build stage
FROM golang:1.25.0-alpine AS builder

WORKDIR /app

# Install Air for live reload (development)
RUN go install github.com/air-verse/air@latest

# Copy go mod and sum
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build the binary (production)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /server ./cmd/server

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the server binary
COPY --from=builder /server .

# Copy Air binary from builder
COPY --from=builder /go/bin/air /usr/local/bin/air

# Copy migration files (if needed)
COPY migrations ./migrations

EXPOSE 8080

# Default command (production)
CMD ["./server"]