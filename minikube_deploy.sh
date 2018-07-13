#!/usr/bin/env bash

# Set docker env so local docker image can be loaded in pods
eval $(minikube docker-env)

docker build -t gambarini/grpc-demo .

kubectl create clusterrolebinding gambarini-admin-binding --clusterrole=cluster-admin --user=gambarini@gmail.com

kubectl apply -f minikube/rabbitmq.yaml

kubectl apply -f minikube/mongodb.yaml

kubectl apply -f minikube/message.yaml

kubectl apply -f minikube/chat.yaml

kubectl apply -f minikube/ingress.yaml