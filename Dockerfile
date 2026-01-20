# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code (including docs folder for Swagger)
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/.env* ./
COPY --from=builder /app/docs ./docs

# Create uploads directory
RUN mkdir -p ./uploads

EXPOSE 8080

CMD ["./main"]