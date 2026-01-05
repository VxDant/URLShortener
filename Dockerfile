FROM golang:1.25-alpine AS builder

# Set environment variables
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set working directory inside the container
WORKDIR /build

# Copy go.mod file for dependency installation
COPY go.mod  ./ go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire application source
COPY . .

# Build the Go binary
RUN go build -o /app .

# Final lightweight stage
FROM alpine:3.21 AS final

# Create non-root user
RUN addgroup -g 1001 app && \
    adduser -D -u 1001 -G app app

# Copy the compiled binary from the builder stage
COPY --from=builder /app /bin/app
COPY --from=builder /build/db/migrations /db/migrations

USER app

# Expose the application's port
EXPOSE 8080

# Run the application
CMD ["/bin/app"]