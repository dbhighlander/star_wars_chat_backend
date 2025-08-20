FROM golang:1.24-alpine AS builder
RUN apk add --no-cache git

WORKDIR /app

COPY . .
RUN ls /app
RUN go mod download



RUN go build -o app .

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/app .
EXPOSE 8080
CMD ["./app"]