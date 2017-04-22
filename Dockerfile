FROM golang:alpine

MAINTAINER sine "sinerwr@gmail.com"

RUN apk --update add git
RUN go-wrapper download github.com/SiCo-DevOps/He
RUN apk del git

WORKDIR $GOPATH/src/github.com/SiCo-DevOps/He

RUN go-wrapper install

WORKDIR $GOPATH/bin/

RUN rm -rf $GOPATH/src

ADD config.sample.json $GOPATH/bin/config.json

EXPOSE 6666

VOLUME $GOPATH/bin/config.json

CMD He