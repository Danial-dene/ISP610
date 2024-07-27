# Use the official Golang image to create a build artifact.
FROM golang:1.20 as builder

# Create and change to the app directory.
WORKDIR /app

# Retrieve application dependencies.
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy local code to the container image.
COPY . ./

# Build the binary for the application named tender-scraper.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o tender-scraper

# Use a slim Debian image for the runtime.
FROM debian:buster-slim
WORKDIR /app

# Copy the compiled application from the builder stage.
COPY --from=builder /app/tender-scraper /app/tender-scraper
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy .env file into the container image
COPY .env /app/

# Run the application
CMD ["./tender-scraper"]
