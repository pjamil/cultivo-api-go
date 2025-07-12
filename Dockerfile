# syntax=docker/dockerfile:1.4

# Stage 1: Builder
FROM golang:1.24.2-alpine AS builder

# Set necessary environment variables for Go modules
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./ 

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the application
# Use -ldflags "-s -w" to strip debug information and symbol table for smaller binary size
# Use -a -installsuffix cgo to force rebuilding packages from source
RUN go build -ldflags "-s -w" -o /cultivo-api-go ./cmd/cultivo-api-go

# Install golang-migrate
RUN go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Stage 2: Runner
FROM alpine:latest

# Install ca-certificates for HTTPS connections
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /cultivo-api-go .

# Expose the port the application listens on
EXPOSE 8080

# Command to run the executable
ENTRYPOINT ["/app/cultivo-api-go"]