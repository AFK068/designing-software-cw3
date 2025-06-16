FROM golang:1.24 as builder

ENV GOPATH=/
WORKDIR /app
COPY ./ /app

RUN go mod download && CGO_ENABLED=0 go build -o /payment ./cmd/run/main.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder /payment /app/payment
COPY ./config/dev.yaml /app/config/dev.yaml
COPY ./db/migrations /app/db/migrations

CMD ["/app/payment"]