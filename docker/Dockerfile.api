FROM golang:1.22-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git ca-certificates
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/bin/api ./cmd/api

FROM alpine:3.19
WORKDIR /app
RUN addgroup -S chengetai && adduser -S chengetai -G chengetai && apk add --no-cache ca-certificates
COPY --from=builder /app/bin/api /usr/local/bin/api
USER chengetai
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/api"]
