# FROM golang:1.24-alpine AS dev

# RUN apk add --no-cache git
# RUN go install github.com/air-verse/air@latest

# WORKDIR /app

# # Copy all source code (assumes go.mod & go.sum are in chat-service/)
# COPY . .

# # Download dependencies
# RUN go mod tidy
# RUN go build -o server ./...

# # Expose your app port
# EXPOSE 8080

# # Run Air for hot reload in dev
# CMD ["air"]

# Build stage
FROM golang:1.24-alpine AS builder

# Install git for go modules
RUN apk add --no-cache git

WORKDIR /app

# Copy all source code (assumes go.mod & go.sum are in chat-service/)
COPY . .

# Download dependencies
RUN go mod tidy
RUN go build -o server ./...

# Final stage: minimal image
FROM alpine:3.18

WORKDIR /app

# Copy the built Go binary
COPY --from=builder /app/server .

# Expose your app port
EXPOSE 8080

# Run the server
CMD ["./server"]