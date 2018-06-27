# GRPC Demo

Proof of concept demo for a simple message application.

It's a microservice architecture, with GRPC, deployed on kubernetes.

For the demo, messages are stored in a MongoDB cluster.

The message routing is done throught GRPC streams connected to RabbitMQ fanout exchanges. Allowing
scalling and availability by a RabbitMQ cluster, with replication of exchanges.


## Deploying

The backend is setup to run on minikube or GKE (Google Kubernetes Engine).
Be sure you 'kubectl'  is configured to the context you need.

Also you will need the protoc cmd to generate the protobuffer GRPC code.

First you have to deploy the NGINX ingress controller


### NGINX Ingress Controller

Cluster Permission for your user (replace marks <> with you data)

```
    kubectl create clusterrolebinding <your-user-cluster-admin-binding> --clusterrole=cluster-admin --user=<your.google.cloud.email@example.org>
```

Mandatory resources for NGINX Ingress Controller (https://github.com/kubernetes/ingress-nginx/blob/master/docs/deploy/index.md)

```
    kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/mandatory.yaml
```

NGINX GKE Ingress LoadBalancer Service

```
    kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/provider/cloud-generic.yaml
```

Or for Minikube:

```
    kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/provider/baremetal/service-nodeport.yaml
```

NGINX nginx-configuration ConfigMap

```
    data:
      client-body-timeout: "3600"
      client-header-timeout: "3600"
      error-log-level: info
      keep-alive: "3600"
      proxy-buffering: "off"
      proxy-request-buffering: "off"
      proxy-stream-timeout: "3600"
      ssl-session-timeout: "3600"
```

### Minikube


```
    $ ./minikube_deploy.sh
```


### GKE

```
    $ ./GKE_deploy.sh
```


## Command line Client

With the backend deployed, you can start the command line client, and start chatting:

```
    $ cd cmdclient
    $ go run main.go
```


## Services

The application is composed of:

- Chat service: The routing hub for all clients. It manages a RabbitMQ fanout exchange
for each client connection.

- Message service: Stores messages and let clients fetch messages by filter.

- Contact service: Not used yet. Will manage contacts.

## Client cmd

The client cmd let you create a contact based on ID. Each contact ID
is unique in the application. You can't send messages to contacts that haven't
been created yet.

Once you connect with one ID, you can reconnect with that same contact in
the future to check messages from other contacts.

You must know a contact ID to send messages to that contact. There is no
contact discovery yet (it's part of the contact service).

## TODO and Road map

- Better Logging, with request correlationID
- Add Unit Testes and Integration Testes
- Slimmer Docker image for services (just copy the binary and run it)
- Solve NGINX ingress 1 min connection timeout
- Add contact discovery and mangement.
- JWT token authentication
- A fancy web client application.

