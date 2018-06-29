FROM golang:alpine

ENV PROJECT /go/src/github.com/zcong1993/ticker

COPY . $PROJECT

WORKDIR ${PROJECT}

RUN go build -o ticker main.go
CMD ["./ticker"]
