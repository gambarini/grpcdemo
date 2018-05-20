#!/usr/bin/env bash

echo Installing gen go
go get -u github.com/golang/protobuf/protoc-gen-go

echo Generating chat golang proto files
protoc ./pb/chat/chat.proto -I. --go_out=plugins=grpc:$GOPATH/src

echo Generating contact golang proto files
protoc ./pb/contact/contact.proto -I. --go_out=plugins=grpc:$GOPATH/src

#echo Generating chat javascript proto files
#mkdir -p js_client
#protoc ./proto/chat/chat.proto -I. --js_out=./js_client --plugin=grpc