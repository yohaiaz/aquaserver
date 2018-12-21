FROM golang:1.9

WORKDIR $GOPATH/src/aqua

COPY . .

RUN go build
RUN chmod +x aquaserver

CMD ["echo hi "]