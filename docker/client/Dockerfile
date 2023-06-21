FROM golang:alpine

RUN mkdir /app
ADD . /app/
WORKDIR /app

RUN go build ./cmd/client

CMD ["/app/client", "-host", "server"]
