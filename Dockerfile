FROM golang:1.21-alpine AS builder
WORKDIR /app
ENV CGO_ENABLED=0
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /bin/server ./cmd/server

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=builder /bin/server /usr/local/bin/server
COPY migrations /app/migrations
WORKDIR /app
EXPOSE 8080
CMD ["/usr/local/bin/server"]
