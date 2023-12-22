FROM golang:latest as builder

WORKDIR /app

COPY go.mod go.sum .

RUN go mod download

COPY . .

RUN make build

FROM alpine:latest

ENV DATA=/data

EXPOSE 80

WORKDIR /app

COPY --from=builder /app/bin/trakr .

ENTRYPOINT ["/app/trakr"]
