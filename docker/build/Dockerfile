FROM golang:1.8.1-alpine

MAINTAINER Cheng-Lung Sung <clsung@gmail.com>

RUN apk add --update git
RUN go get -u github.com/csigo/authsrv/cmd/csigosrv

CMD ${GOPATH}/bin/csigosrv
