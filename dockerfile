# Use an official Go runtime as a parent image
FROM golang:1.22.3-alpine

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go applications
RUN go build -o broker ./cmd/broker
RUN go build -o consumer ./cmd/consumer
RUN go build -o producer ./cmd/producer

# Expose port 8080 to the outside world
EXPOSE 8080

# Default command
CMD ["./broker"]

