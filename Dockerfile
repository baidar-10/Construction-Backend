# Build stage
FROM golang:1.24-bookworm AS builder

WORKDIR /app

# Install build dependencies
RUN apt-get update \
	&& apt-get install -y --no-install-recommends git ca-certificates \
	&& rm -rf /var/lib/apt/lists/*

# Work around MITM/unknown CA by bypassing proxy/sumdb
ENV GOPROXY=direct \
	GOSUMDB=off \
	GONOPROXY=* \
	GONOSUMDB=* \
	GOPRIVATE=* \
	GOINSECURE=* \
	GIT_SSL_NO_VERIFY=1

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code (including docs folder for Swagger)
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Final stage
FROM debian:bookworm-slim

RUN apt-get update \
	&& apt-get install -y --no-install-recommends ca-certificates tzdata \
	&& rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/.env* ./
COPY --from=builder /app/docs ./docs

# Create uploads directory
RUN mkdir -p ./uploads

EXPOSE 8080

CMD ["./main"]