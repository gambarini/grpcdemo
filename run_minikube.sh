#!/usr/bin/env bash

./genproto.sh

# Set docker env so local docker image can be loaded in pods
eval $(minikube docker-env)

docker build -t gambarini/grpc-demo .

kubectl apply -f minikube/chat.yaml
kubectl apply -f minikube/contact.yaml

chat_url="$(minikube service chat-service --url)"
echo chat-service on: ${chat_url}

chat_url="$(minikube service contact-service --url)"
echo chat-contact on: ${chat_url}

read -p "Press any key to stop serving... " -n1 -s

kubectl delete deployment chat
kubectl delete service chat-service
kubectl delete deployment contact
kubectl delete service contact-service
