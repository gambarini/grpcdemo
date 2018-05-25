#!/usr/bin/env bash

kubectl create namespace rabbitmq

kubectl run etcd --image=microbox/etcd --port=4001 --namespace=rabbitmq -- --name etcd
kubectl --namespace=rabbitmq expose deployment etcd

eval $(minikube docker-env)
docker build . -t rabbitmq-autocluster

kubectl create secret generic --namespace=rabbitmq erlang.cookie --from-file=./erlang.cookie

kubectl create -f rabbitmq.yaml

