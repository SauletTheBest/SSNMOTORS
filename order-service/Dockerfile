# Start from the official Go image
FROM golang:1.21-alpine

# Set the working directory inside the container
WORKDIR /app

# Install Git and build tools
RUN apk update && apk add --no-cache git

# Copy go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire application source code
COPY . .

# Build the Go application
RUN go build -o order-service ./cmd

# Expose the gRPC port
EXPOSE 50051

# Command to run the executable
CMD ["./order-service"]
