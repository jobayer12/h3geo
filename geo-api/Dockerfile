# Multi-stage build for Angular and Go
FROM node:18-alpine AS angular-builder

# Set working directory for Angular build
WORKDIR /app/geo-map-app

# Copy package files
COPY geo-map-app/package*.json ./

# Install dependencies
RUN npm ci

# Copy Angular source code
COPY geo-map-app/ ./

# Build Angular app
RUN npm run build

# Go build stage
FROM golang:1.22-alpine AS go-builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy Go module files
COPY geo-api/go.mod geo-api/go.sum ./

# Download dependencies
RUN go mod download

# Copy Go source code
COPY geo-api/ ./

# Copy built Angular files from previous stage
COPY --from=angular-builder /app/geo-map-app/dist/geo-map-app/browser ./static/browser

# Build Go application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o geo-api .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create app user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy the binary from go-builder stage
COPY --from=go-builder /app/geo-api .

# Copy static files from go-builder stage
COPY --from=go-builder /app/static ./static

# Change ownership to app user
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Run the application
CMD ["./geo-api"]
