#!/usr/bin/env bash

kubectl create namespace rabbitmq

kubectl run etcd --image=microbox/etcd --port=4001 --namespace=rabbitmq -- --name etcd
kubectl --namespace=rabbitmq expose deployment etcd

eval $(minikube docker-env)
docker build . -t rabbitmq-autocluster

kubectl create secret generic --namespace=rabbitmq erlang.cookie --from-file=./erlang.cookie

kubectl create -f rabbitmq.yaml


FIRST_POD=$(kubectl get pods --namespace rabbitmq -l 'app=rabbitmq' -o jsonpath='{.items[0].metadata.name }')
kubectl exec --namespace=rabbitmq $FIRST_POD rabbitmqctl cluster_status