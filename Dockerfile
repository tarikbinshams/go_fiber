# Use official Go image to build the Go app
FROM golang:1.20-alpine as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod tidy

# Copy the entire project
COPY . .

# Build the Go application
RUN go build -o main .

# Start a new stage from a clean Alpine image to run the app
FROM alpine:latest

# Install necessary dependencies for the app (e.g., curl, bash)
RUN apk --no-cache add ca-certificates

# Copy the built binary from the builder image
COPY --from=builder /app/main /app/

# Set environment variables for the application (e.g., MongoDB URI)
ENV MONGO_URI=mongodb://mongo:27017

# Expose the port the app will run on
EXPOSE 3000

# Run the Go app
CMD ["/app/main"]
