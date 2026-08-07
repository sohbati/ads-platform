FROM golang:1.25-alpine AS builder

WORKDIR /src
COPY nats-message-broker/go.mod nats-message-broker/go.sum ./
RUN go mod download
COPY nats-message-broker/ ./
RUN CGO_ENABLED=0 go build -o /server ./cmd/server

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /server ./server
EXPOSE 8095
ENV PORT=8095
ENV NATS_HOST=0.0.0.0
ENV NATS_PORT=4222
ENV NATS_HTTP_PORT=8222
CMD ["./server"]
