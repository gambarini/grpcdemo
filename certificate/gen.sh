#!/usr/bin/env bash

rm ca.*
rm server.*


openssl genrsa -out ca.key 2048

openssl req -x509 -new -nodes -key ca.key -subj "/CN=grpcdemo.com" -days 10000 -out ca.crt

openssl genrsa -out server.key 2048

openssl req -new -key server.key -out server.csr -config csr.conf

openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 10000 -extensions v3_ext -extfile csr.conf

openssl x509  -noout -text -in ./server.crt

kubectl delete secret tls.ingress.secret
kubectl create secret tls tls.ingress.secret --key server.key --cert server.crt

kubectl delete secret tls.authority.secret
kubectl create secret tls tls.authority.secret --key ca.key --cert ca.crt


