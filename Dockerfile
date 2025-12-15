FROM golang:1.25.4-alpine AS builder

WORKDIR /app

# Copy entire project structure FIRST
COPY . .

# Then run go mod download (after all files are copied)
RUN go mod download

# Build the application
RUN go build -o app

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/app .
CMD ["./app"]