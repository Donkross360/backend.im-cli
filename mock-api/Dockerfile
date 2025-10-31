FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files first
COPY go.mod ./

# Download dependencies (creates go.sum if needed)
RUN go mod download

# Copy source code
COPY . .

# Ensure go.sum is up to date and build
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -o mock-api main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/mock-api .

EXPOSE 8080

CMD ["./mock-api"]

