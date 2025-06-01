# Build stage
FROM golang:1.24.3-alpine3.22 AS builder

WORKDIR /app

# Install dependencies including swag
RUN apk add --no-cache git make gcc musl-dev && \
    go install github.com/swaggo/swag/cmd/swag@latest

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN swag init -g cmd/cultivo-api-go-swagger/main.go && \
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o cultivo-api-go ./cmd/cultivo-api-go

# Runtime stage
FROM alpine:3.22.0

WORKDIR /root/

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/cultivo-api-go .
COPY --from=builder /app/.env-prod .

# Expose port
EXPOSE 8080

# Command to run the executable
CMD ["./cultivo-api-go"]
