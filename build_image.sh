#!/usr/bin/env bash

./genproto.sh

docker build -t gambarini/grpc-demo .