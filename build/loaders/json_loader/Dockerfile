FROM golang:1.23.2-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY pkg/ ./pkg/
COPY cmd/loaders/json_loader/ ./cmd/loaders/json_loader/
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o loader ./cmd/loaders/json_loader

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/loader ./
COPY assets/data/big.txt assets/data/big.txt
CMD ["./loader"]
