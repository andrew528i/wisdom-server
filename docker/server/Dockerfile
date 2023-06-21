FROM golang:alpine

RUN mkdir /app
ADD . /app/
WORKDIR /app

RUN go build ./cmd/server

CMD ["/app/server"]
