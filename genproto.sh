#!/usr/bin/env bash

echo Generating chat golang proto files
mkdir -p chat/pb
protoc ./proto/chat.proto -I. --go_out=plugins=grpc:$GOPATH/src

echo Generating contact golang proto files
mkdir -p contact/pb
protoc ./proto/contact.proto -I. --go_out=plugins=grpc:$GOPATH/src

#echo Generating chat javascript proto files
#mkdir -p js_client
#protoc ./proto/chat/chat.proto -I. --js_out=./js_client --plugin=grpc