FROM golang:1.24-alpine AS dev

RUN apk add --no-cache git
RUN go install github.com/air-verse/air@latest

WORKDIR /app

# Copy go.mod and go.sum for dependency caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of your source code
COPY . .

# Expose your app port
EXPOSE 8080

# Run Air for hot reload in dev
CMD ["air"]
