FROM golang:1.9-alpine

ENV GOPATH /go

WORKDIR /go/src/github.com/gambarini/grpcdemo

COPY . .


RUN apk add --no-cache git && go get -d -v ./...


EXPOSE 80

ENTRYPOINT ["go", "run"]


