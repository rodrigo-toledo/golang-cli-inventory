FROM --platform=linux/arm64 golang:1.25 AS builder

WORKDIR /app

# Set the GOEXPERIMENT environment variable for JSON v2
ENV GOEXPERIMENT=jsonv2

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o cli-inventory cmd/inventory/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/cli-inventory .

# Command to run the executable
CMD ["./cli-inventory"]
