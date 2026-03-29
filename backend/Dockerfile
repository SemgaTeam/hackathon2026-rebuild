FROM golang:1.26.1-alpine AS builder
WORKDIR /usr/src/app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o application ./cmd/

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /usr/src/app/application ./app

RUN chown -R appuser:appgroup /app

USER appuser

EXPOSE 8080

CMD ["./app"]
