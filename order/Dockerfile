FROM golang:1.24 as builder

ENV GOPATH=/
WORKDIR /app
COPY ./ /app

RUN go mod download && CGO_ENABLED=0 go build -o /order ./cmd/run/main.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder /order /app/order
COPY ./config/dev.yaml /app/config/dev.yaml
COPY ./db/migrations /app/db/migrations

CMD ["/app/order"]