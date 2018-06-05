# GRPC Demo

Proof of concept demo for a simple message application.

It's a microservice architecture, with GRPC, deployed on kubernetes.

For the demo, messages are stored in a MongoDB cluster.

The message routing is done throught GRPC streams connected to RabbitMQ fanout exchanges. Allowing
scalling and availability by a RabbitMQ cluster, with replication of exchanges.


## Running

The backend is setup to run on minikube. So you will need a minikube kube cluster
running. It's recomended to set the memory to 4096 (minikube start --memory 4096).

Also you will need the protoc cmd to generate the protobuffer GRPC code.

Once minikube is on, deploy the RabbitMQ cluster:

```
    $ cd rabbitmq
    $ ./run_rabbitmq.sh
```

Then deploy the mongoDB cluster:

```
    $ cd ..
    $ ./run_mongodb.sh
```

Finally deploy the services:

```
    $ ./run_services
```

Now you can start the command line client, and start chatting:

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

You must to know a contact ID to send messages to that contact. There is no
contact discovery yet (it's part of the contact service).

## Future changes

- Add contact Discovery and mangement
- Create a Google Kubernetes Engine deployment
