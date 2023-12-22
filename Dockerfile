FROM golang:latest as builder

WORKDIR /app

COPY go.mod go.sum .

RUN go mod download

COPY . .

RUN make build

FROM debian:stable-slim

ENV DATA=/data

EXPOSE 80

RUN mkdir /data

WORKDIR /app

COPY --from=builder /app/bin bin
COPY --from=builder /app/public public

RUN chmod +x /app/bin/trakr


ENTRYPOINT ["/app/bin/trakr"]
