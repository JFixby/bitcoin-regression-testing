FROM golang:1.11

WORKDIR /go/src/github.com/jfixby/btcregtest
COPY . .

RUN apt-get update && apt-get upgrade -y && apt-get install -y rsync

RUN git clone https://github.com/btcsuite/btcd /go/src/github.com/btcsuite/btcd

RUN cd /go/src/github.com/btcsuite/btcd && env GO111MODULE=on go install . .\cmd\...

RUN btcd --version
