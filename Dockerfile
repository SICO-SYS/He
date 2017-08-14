FROM golang:alpine

MAINTAINER sine "sinerwr@gmail.com"

RUN apk --update add git && \
    go-wrapper download github.com/SiCo-Ops/He && \
    apk del git && \
    cd $GOPATH/src/github.com/SiCo-Ops/He && \
    go-wrapper install && \
    rm -rf $GOPATH/src

EXPOSE 6666

VOLUME $GOPATH/bin/config.json

CMD He