FROM golang:1.4
RUN mkdir -p /go/src/app
WORKDIR /go/src/app

ADD captain /go/src/app/captain/
ADD main.go /go/src/app/

RUN go-wrapper download
RUN go-wrapper install
