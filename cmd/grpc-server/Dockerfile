FROM golang:1.24.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/grpc-server ./cmd/grpc-server
COPY gen ./gen
COPY proto ./proto

RUN go build -o grpc-server ./cmd/grpc-server

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/grpc-server .

EXPOSE 50051

CMD ["./grpc-server"]