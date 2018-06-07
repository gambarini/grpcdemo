#!/usr/bin/env bash

kubectl run etcd --image=microbox/etcd --port=4001 -- --name etcd
kubectl expose deployment etcd

cd rabbitmq
docker build . -t gambarini/rabbitmq-autocluster
docker push gambarini/rabbitmq-autocluster
cd ..

kubectl create secret generic erlang.cookie --from-file=rabbitmq/erlang.cookie

kubectl create -f GKE/rabbitmq.yaml

kubectl apply -f GKE/mongodb.yaml

docker build -t gambarini/grpc-demo .
docker push gambarini/grpc-demo

kubectl apply -f GKE/chat.yaml
kubectl apply -f GKE/message.yaml

kubectl apply -f GKE/chatLB.yaml
kubectl apply -f GKE/messageLB.yaml