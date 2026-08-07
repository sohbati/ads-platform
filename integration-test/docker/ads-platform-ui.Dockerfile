FROM golang:1.24-alpine AS builder

WORKDIR /src
COPY ads-platform-ui/go.mod ads-platform-ui/go.sum ./
RUN go mod download
COPY ads-platform-ui/ ./
RUN CGO_ENABLED=0 go build -o /server ./cmd/server

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /src ./
COPY --from=builder /server ./server
EXPOSE 8094
ENV PORT=8094
CMD ["./server"]
