# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Install migrate tool
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./server

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates postgresql-client

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/server .
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/.env .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./server"]
