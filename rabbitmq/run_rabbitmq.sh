#!/usr/bin/env bash


kubectl run etcd --image=microbox/etcd --port=4001 -- --name etcd
kubectl expose deployment etcd

eval $(minikube docker-env)
docker build . -t rabbitmq-autocluster

kubectl create secret generic erlang.cookie --from-file=./erlang.cookie

kubectl create -f rabbitmq.yaml

