FROM golang:1.23.1-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/api ./cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/bin/api /app/api
COPY --from=builder /app/migrations /app/migrations
COPY ./local.env /app/local.env

EXPOSE 3000
CMD ["/app/api"]