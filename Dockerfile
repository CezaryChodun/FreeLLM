FROM golang:1.26-alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /freellm ./cmd/server

FROM alpine:3.21

RUN apk add --no-cache ca-certificates
COPY --from=builder /freellm /usr/local/bin/freellm
COPY defaults/ /app/defaults/
COPY config.yml /app/config.yml

WORKDIR /app

EXPOSE 3000

ENTRYPOINT ["freellm"]
