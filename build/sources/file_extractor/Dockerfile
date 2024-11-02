FROM golang:1.23.2-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY pkg/ ./pkg/
COPY cmd/sources/file_extractor/ ./cmd/sources/file_extractor/
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o extractor ./cmd/sources/file_extractor

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/extractor ./
COPY assets/data/big.txt assets/data/big.txt
CMD ["./extractor"]
