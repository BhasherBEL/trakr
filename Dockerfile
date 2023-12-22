FROM golang:latest as builder

WORKDIR /app

COPY go.mod go.sum .

RUN go mod download

COPY . .

RUN make build

FROM debian:stable-slim

ENV DATA=/data

EXPOSE 80

WORKDIR /app

COPY --from=builder /app/bin/trakr .

RUN chmod +x /app/trakr

RUN mkdir /data

ENTRYPOINT ["/app/trakr"]
