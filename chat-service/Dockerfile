FROM golang:1.24-alpine AS builder
RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app .

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/app .
EXPOSE 8080
CMD ["./app"]