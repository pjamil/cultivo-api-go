# Build stage
FROM golang:1.24.3-alpine3.22 AS builder

WORKDIR /app

# Copy the source code
COPY . .

# Install dependencies including swag
RUN apk add --no-cache gcc git make musl-dev && \
    go install github.com/swaggo/swag/cmd/swag@latest

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Build the application
RUN swag init -g cmd/cultivo-api-go/main.go && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o cultivo-api-go ./cmd/cultivo-api-go

# Runtime stage
FROM alpine:3.22.0

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app

# Set the timezone to Sao Paulo
ENV TZ=America/Sao_Paulo

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/cultivo-api-go .
COPY --from=builder /app/.env .

HEALTHCHECK --interval=30s --timeout=3s \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Expose port
EXPOSE 8080

# Command to run the executable
CMD ["./cultivo-api-go"]
