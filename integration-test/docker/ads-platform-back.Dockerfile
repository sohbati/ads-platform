FROM golang:1.25-alpine AS builder

WORKDIR /src
COPY ads-platform-back/go.mod ads-platform-back/go.sum ./
RUN go mod download
COPY ads-platform-back/ ./
RUN CGO_ENABLED=0 go build -o /server ./cmd/server

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /server ./server
EXPOSE 8092
ENV APPLICATION_SERVER_PORT=8092
CMD ["./server"]
