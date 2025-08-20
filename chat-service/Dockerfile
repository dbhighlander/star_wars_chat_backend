FROM golang:1.24-alpine AS builder
RUN apk add --no-cache git

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Force static build
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

# Check what was built
RUN ls -l /app

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/app .
RUN ls -l /app
CMD ["./app"]