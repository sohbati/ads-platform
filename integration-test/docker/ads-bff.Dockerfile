FROM golang:1.24-alpine AS builder

WORKDIR /src
COPY ads-bff/go.mod ads-bff/go.sum ./
RUN go mod download
COPY ads-bff/ ./
RUN CGO_ENABLED=0 go build -o /server ./cmd/server

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /server ./server
EXPOSE 8097
ENV PORT=8097
CMD ["./server"]
