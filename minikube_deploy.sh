#!/usr/bin/env bash

# Set docker env so local docker image can be loaded in pods
eval $(minikube docker-env)

docker build -t gambarini/grpc-demo .

cd rabbitmq
docker build . -t gambarini/rabbitmq-autocluster
cd ..

kubectl create clusterrolebinding gambarini-admin-binding --clusterrole=cluster-admin --user=gambarini@gmail.com

kubectl apply -f GKE/rabbitmq.yaml

kubectl apply -f GKE/mongodb.yaml

kubectl apply -f GKE/message.yaml

kubectl apply -f GKE/chat.yaml

kubectl apply -f GKE/ingress.yaml