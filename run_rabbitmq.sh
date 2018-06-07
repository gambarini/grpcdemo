#!/usr/bin/env bash


kubectl run etcd --image=microbox/etcd --port=4001 -- --name etcd
kubectl expose deployment etcd

eval $(minikube docker-env)
cd rabbitmq
docker build . -t rabbitmq-autocluster
cd ..

kubectl create secret generic erlang.cookie --from-file=rabbitmq/erlang.cookie

kubectl create -f minikube/rabbitmq.yaml

