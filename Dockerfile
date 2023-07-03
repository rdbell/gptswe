# Use the official Golang image as a base
FROM golang:1.20

# Set the working directory
WORKDIR /app

# Copy the go.mod and go.sum files into the working directory
COPY go.mod .
COPY go.sum .

# Install the dependencies
RUN go mod download

# Copy the remaining source code into the working directory
COPY *.go .

# Build the Go application
RUN go build -o gptswe .

# Command to run the application when the container starts
ENTRYPOINT ["/app/gptswe"]
