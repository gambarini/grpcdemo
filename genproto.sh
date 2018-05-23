#!/usr/bin/env bash

echo Installing gen go
go get -u github.com/golang/protobuf/protoc-gen-go

echo Generating chat golang proto files
protoc ./pb/proto/chat.proto -I. --go_out=plugins=grpc:$GOPATH/src

echo Generating contact golang proto files
protoc ./pb/proto/contact.proto -I. --go_out=plugins=grpc:$GOPATH/src

#echo Generating chatsvc javascript proto files
#mkdir -p js_client
#protoc ./proto/chatsvc/chatsvc.proto -I. --js_out=./js_client --plugin=grpc