FROM golang:1.25-alpine AS builder

WORKDIR /src
COPY ads-platform-notification/go.mod ads-platform-notification/go.sum ./
RUN go mod download
COPY ads-platform-notification/ ./
RUN CGO_ENABLED=0 go build -o /server ./cmd/server

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /server ./server
EXPOSE 8096
ENV PORT=8096
CMD ["./server"]
