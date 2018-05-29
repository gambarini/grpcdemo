#!/usr/bin/env bash

./genproto.sh

go build ./...

# Set docker env so local docker image can be loaded in pods
eval $(minikube docker-env)

docker build -t gambarini/grpc-demo .

kubectl apply -f minikube/chat.yaml
kubectl apply -f minikube/contact.yaml
kubectl apply -f minikube/message.yaml

chat_url="$(minikube service chat-service --url)"
echo chat-service on: ${chat_url}

chat_url="$(minikube service contact-service --url)"
echo contact-service on: ${chat_url}

chat_url="$(minikube service message-service --url)"
echo message-service on: ${chat_url}

read -p "Press any key to stop serving services..." -n1 -s

kubectl delete deployment chat
kubectl delete service chat-service
kubectl delete deployment contact
kubectl delete service contact-service
kubectl delete deployment message
kubectl delete service message-service