# Use a GoLang base image
FROM golang:1.20

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download Go module dependencies
RUN go mod download

# Copy the app source code to the container
COPY . .

# Build the GoLang app
RUN go build -o app

# Set the environment variables
ENV IP_VERSION=4
ENV LEDGER_SERVICE_HOST=0.0.0.0
ENV LEDGER_SERVICE_PORT=50053

# Expose the port on which Ledger Service listens
EXPOSE 50053

# Start the app when the container is run
CMD ["./app"]