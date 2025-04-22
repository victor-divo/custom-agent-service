FROM golang:1.22-alpine

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy semua file
COPY . .

# Build binary
RUN cd cmd/ && go build -o su /app/custom-agent-service

# Jalankan service
CMD ["./custom-agent-service"]
