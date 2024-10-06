FROM golang:1.23.2-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server ./
COPY assets/data/big.txt assets/data/big.txt
CMD ["./server serve --http=0.0.0.0:8000"]
