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
ENV PAYMENT_API_SERVICE_HOST=0.0.0.0
ENV PAYMENT_API_SERVICE_PORT=8080
ENV MERCHANT_SERVICE_HOST=0.0.0.0
ENV MERCHANT_SERVICE_PORT=50051
ENV RATE_LIMITER_SERVICE_HOST=0.0.0.0
ENV RATE_LIMITER_SERVICE_PORT=50052
ENV LEDGER_SERVICE_HOST=0.0.0.0
ENV LEDGER_SERVICE_PORT=50053
ENV BANK_SIMULATOR_HOST=0.0.0.0
ENV BANK_SIMULATOR_PORT=8000

# Expose the port on which the API listens
EXPOSE 8080

# Start the app when the container is run
CMD ["./app"]