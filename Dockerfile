# Use an official Go runtime as a parent image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Disable Go modules for this build
ENV GO111MODULE=off

# Copy the local Go application code to the container
COPY . .

# Build the Go application inside the container
RUN go build -o server

# Expose the port that your Go application listens on (change 8080 to your application's port)
EXPOSE 8080

# Command to run your Go application
CMD ["./server"]



