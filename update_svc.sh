#!/usr/bin/env bash

./genproto.sh

go build ./...

# Set docker env so local docker image can be loaded in pods
eval $(minikube docker-env)
docker build -t gambarini/grpc-demo .

kubectl delete deployment chat
kubectl delete service chat-service
kubectl delete deployment message
kubectl delete service message-service

kubectl apply -f minikube/chat.yaml
kubectl apply -f minikube/message.yaml

