FROM golang:1.24-alpine AS builder

WORKDIR /src
COPY ads-platform-cache-service/go.mod ads-platform-cache-service/go.sum ./
RUN go mod download
COPY ads-platform-cache-service/ ./
RUN CGO_ENABLED=0 go build -o /server ./cmd/server

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /server ./server
EXPOSE 8093
ENV PORT=8093
CMD ["./server"]
