# Start with the official Golang base image
FROM golang:1.16-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o exoplanet-service

# Expose the port the service will run on
EXPOSE 8080

# Command to run the executable
CMD ["./exoplanet-service"]
